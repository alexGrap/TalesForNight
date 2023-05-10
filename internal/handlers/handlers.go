package handlers

import (
	"fmt"
	"fsm/internal/keyboards"
	"fsm/internal/models"
	"fsm/internal/usecase"
	"fsm/pkg/repository"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"strconv"
)

var (
	BaseSG     = fsm.NewStateGroup("base")
	TaleState  = BaseSG.New("Tale")
	SpeakState = BaseSG.New("Speak")
	GenreState = BaseSG.New("Genre")
	OtherState = BaseSG.New("Other")
	AdminState = BaseSG.New("Admin")
)

func StartHandlers(bot *tele.Group, manager *fsm.Manager) {
	bot.Handle("/start", onStart)
	bot.Handle("/admin", sendAdmin)
	bot.Handle("/help", helper)
	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.FSMContext) error {
		s := state.State()
		return c.Send(s.String())
	})

	// buttons
	manager.Bind(&keyboards.FindBtn, fsm.AnyState, onTaleChoose(keyboards.OurTaleBtn, keyboards.OwnTaleBtn,
		keyboards.CancelBtn))
	manager.Bind(&keyboards.GenreBtn, fsm.AnyState, onGenreChoose(keyboards.FairyBtn, keyboards.PoemBtn,
		keyboards.DramaBtn, keyboards.CancelBtn))
	manager.Bind(&keyboards.SpeechBtn, fsm.AnyState, onSpeechChoose(keyboards.TextBtn, keyboards.YandexBtn,
		keyboards.CancelBtn))
	manager.Bind(&keyboards.InfoBtn, fsm.AnyState, onInfoChoose(keyboards.GetSleepingInfoBtn,
		keyboards.SleepingAdviceBtn, keyboards.CancelBtn))
	manager.Bind(&keyboards.CancelBtn, fsm.AnyState, onCancelForm())
	manager.Bind(&keyboards.UserBtn, fsm.AnyState, userInformation)

	//// form
	manager.Bind(&keyboards.TextBtn, SpeakState, setSpeak("Текст"))
	manager.Bind(&keyboards.YandexBtn, SpeakState, setSpeak("Yandex"))

	manager.Bind(&keyboards.DramaBtn, GenreState, setGenre("Драма"))
	manager.Bind(&keyboards.FairyBtn, GenreState, setGenre("Сказка"))
	manager.Bind(&keyboards.PoemBtn, GenreState, setGenre("Поэма"))

	manager.Bind(&keyboards.OurTaleBtn, TaleState, generateTail)
	manager.Bind(&keyboards.OwnTaleBtn, TaleState, waitOwnState(keyboards.CancelBtn))
	manager.Bind(tele.OnText, TaleState, choosingTitle)

	manager.Bind(&keyboards.GetSleepingInfoBtn, OtherState, sendInfo)
	manager.Bind(&keyboards.SleepingAdviceBtn, OtherState, sendAdvice)

	manager.Bind(&keyboards.AdminYandexBtn, fsm.AnyState, yandexToZero)
	manager.Bind(&keyboards.AdminSendlerBtn, fsm.AnyState, startSendler(keyboards.CancelBtn))
	manager.Bind(tele.OnText, AdminState, sendler)

}

func onStart(c tele.Context) error {
	var body models.User
	body.UserId = c.Sender().ID
	repository.CreateUser(body)
	log.Println("new user", c.Sender().ID)

	c.Send(
		fmt.Sprintf("Добро пожаловать в бот-рассказчик, %s 📕\n", c.Sender().FirstName), keyboards.OnStartKB())

	return c.Send("Данный бот был запущен в рамках тестового задания для VK. Изначальная задумка была создать удобного бота, " +
		"который бы искал подходящее произведение для каждого пользователя, и при желании делал бы из него аудиокнигу. Так " +
		"как бОльшая часть литературы находится под действием авторского права, была использованна технология" +
		" Chat GPT от OpenAi для демонстрации возможностей бота, а точнее для генерации случайного небольшого отрывка." +
		"Для синтеза речи использованa разработка компании " +
		"Яндекс - Yandex SpeechKit. Так как последняя технология предоставляется на коммерческой основе, количество " +
		"пользований данной озвучкой текста ограниченно 15 (для дополнительной тестровки - свяжитесь со мной, я обнулю). Данные два способа были выбраны в демонстрационных целях, " +
		"и всегда могут быть заменены на аналоги, например на api \"Маруси\", доступ к которой предоставляется на той " +
		"же основе, что и у Яндекса. Приятного пользования!\n\n Для ознакомления с функционалом рекомендуем /help")

}

func helper(c tele.Context) error {
	return c.Send("Справка для пользователей 📃\n\nВ разделе \"Выбрать жанр 🖋\" Вы можете выбрать жанр произведения, которое будет сгенерированно" +
		" случайным образом специально для Вас (по умолчанию \"Сказка\")\n\nВ разделе \"Задать озвучку 🎙\" Вы можете выбрать между " +
		"текстом или ограниченной 5ью запросами озвучкой Yandex SpeechKit (по умолчанию \"Текст\")\n\nВ " +
		"разделе \"Дополнительная информация 🕶\" Вы можете прочитать интересные факты о сне и рекомендации о том, как " +
		"быстрее уснуть\n\nВ разделе \"Ваши настройки ⚙\" хранится информация о Ваших текущих установленных " +
		"параметрах\n\nИ наконец, в разделе \"Найти скзаку 🔍\" Вы можете запросить случайно сгенерированный фрагмент" +
		" того жанра, который выбран Вами в разделе \"Выбрать жанр 🖋\", или же запросить сгенерировать отрывок из " +
		"Вашего произведения.")
}

func onTaleChoose(ownBtn tele.Btn, ourBtn tele.Btn, cancel tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(ownBtn, ourBtn), menu.Row(cancel))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(TaleState)
		return c.Send("Выберите, какую произведение Вы прослушаете сегодня", menu)
	}
}

func onGenreChoose(tale tele.Btn, poem tele.Btn, drama tele.Btn, cancel tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(tale, poem),
		menu.Row(drama, cancel))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(GenreState)
		return c.Send("Выберите жанр Вашего произведения:", menu)

	}
}

func onSpeechChoose(python tele.Btn, yandex tele.Btn, cancel tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(python, yandex),
		menu.Row(cancel))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(SpeakState)
		return c.Send("Выберите озвучку Вашего произведения:", menu)

	}
}

func onInfoChoose(info tele.Btn, advice tele.Btn, cancel tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(info, advice), menu.Row(cancel))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(OtherState)
		return c.Send("Выберите, какую информацию Вам было бы интересно прочитать", menu)

	}
}

func onCancelForm() fsm.Handler {
	menu := keyboards.OnStartKB()
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(OtherState)
		return c.Send("Okay", menu)

	}
}

func setSpeak(speaker string) fsm.Handler {
	menu := keyboards.OnStartKB()
	return func(c tele.Context, state fsm.FSMContext) error {
		repository.UpdateSounder(c.Sender().ID, speaker)
		return c.Send(fmt.Sprintf("Установлен формат: %s", speaker), menu)
	}
}

func setGenre(genre string) fsm.Handler {
	menu := keyboards.OnStartKB()
	return func(c tele.Context, state fsm.FSMContext) error {
		repository.UpdateGenre(c.Sender().ID, genre)
		return c.Send(fmt.Sprintf("Установлен жанр %s", genre), menu)
	}
}

func choosingTitle(c tele.Context, state fsm.FSMContext) error {
	title := c.Message().Text
	menu := keyboards.OnStartKB()
	c.Send("Просим прощения, генерация текста и аудио занимает время. Пожалуйста, подождите")
	body := repository.GetUser(c.Sender().ID)
	repository.UpdateBook(c.Sender().ID, title)
	message := usecase.GenerateTale(fmt.Sprintf("Прочитай отрывок в 10 абзацев из %s", title), body)
	if message == "." {
		fileSendler(c, state)
		return nil
	}
	return c.Send(fmt.Sprintf("Мы постараемся найти для Вас \"%s.\n\n%s", title, message), menu)
}

func waitOwnState(cancel tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancel))
	return func(c tele.Context, state fsm.FSMContext) error {
		return c.Send("Введите название книги, которая Вам интересна", menu)
	}
}

func fileSendler(c tele.Context, state fsm.FSMContext) error {
	menu := keyboards.OnStartKB()

	fileCloser(c, state)
	path, _ := os.Getwd()
	err := os.Remove(path + "/temp-folder/file.ogg")
	if err != nil {
		log.Println(err)
		return c.Send("Простите, но у нас произошла ошибка. Мы обещаем исправиться)", menu)
	}
	return nil
}

func generateTail(c tele.Context, state fsm.FSMContext) error {
	menu := keyboards.OnStartKB()
	body := repository.GetUser(c.Sender().ID)
	c.Send("Просим прощения, генерация текста и аудио занимает время. Пожалуйста, подождите")
	message := usecase.GenerateTale(fmt.Sprintf("Прочитай слуайную %s", body.Genre), body)
	if message == "." {
		fileSendler(c, state)
		return nil
	}
	return c.Send(fmt.Sprintf("Мы подобрали для Вас это: \n\n%s", message), menu)
}

func fileCloser(c tele.Context, state fsm.FSMContext) error {
	path, _ := os.Getwd()
	menu := keyboards.OnStartKB()
	path += "/temp-folder/file.ogg"
	a := &tele.Audio{File: tele.FromDisk(path)}
	log.Println(a.OnDisk())
	if a.OnDisk() {
		c.Send("Приятного прослушивания)")
	}
	return c.Send(a, menu)

}

func userInformation(c tele.Context, state fsm.FSMContext) error {
	body := repository.GetUser(c.Sender().ID)
	menu := keyboards.OnStartKB()
	return c.Send(fmt.Sprintf("Информация о Вас: \n\nВаше имя 👦🏻: %s\nВыбранный жанр 🎭: %s\nВыбранная озвучка"+
		" 🔊: %s\nВыбранная книга 📚: %s\nКоличество использований Yandex: %d/5",
		c.Sender().FirstName, body.Genre, body.Sounder, body.Book, body.Counter), menu)
}

func sendInfo(c tele.Context, state fsm.FSMContext) error {
	menu := keyboards.OnStartKB()
	c.Send("Вот небольшая информация о человеческом сне🌘")
	message := "1. 12% людей видят сны исключительно в черно-белых тонах, в то время как до появления цветного телевидения" +
		" только 15% людей видели сны в цвете.\n" +
		"2. Люди спят 1/3 своей жизни. Очевидно, это зависит от возраста человека, но в среднем составляет около трети," +
		" что довольно много, если подумать.\n 3. Самый продолжительный период без сна - 11 дней. Это было установлено " +
		"калифорнийским студентом по имени Рэнди Гарднер в 1964 году. Не повторяйте это в домашних условиях.\n" +
		"4. Нередко глухие люди используют язык жестов во сне. Есть много случаев, когда люди сообщали о своих" +
		" глухих партнерах или детях, использующих язык жестов во сне.\n" +
		"5. Дисания — состояние, когда утром трудно вставать с кровати. Мы все, несомненно, время" +
		" от времени хотим дольше поспасть, но тем, кто страдает от дисании, это особенно трудно." +
		" Скорее всего, это форма синдрома хронической усталости.\n" +
		"6. Парасомния — термин, обозначающий неестественные движения во время сна. Некоторые люди даже" +
		" совершали преступление из-за парасомнии, включая вождение во сне и даже убийство.\n" +
		"7. Считается, что до 15% населения — лунатики. Существует мнение, что нельзя будить кого-то," +
		" кто ходит во сне, но это не более, чем миф. 10. Каждая четвертая супружеская пара спит в разных кроватях.\n" +
		"8. Лишение сна убивает быстрее, чем лишение пищи.\n" +
		"9. Те, кто родился слепым, переживают сны, связанные с такими вещами, как эмоции, звук, запах, а не зрение."
	return c.Send(message, menu)
}

func sendAdvice(c tele.Context, state fsm.FSMContext) error {
	menu := keyboards.OnStartKB()
	c.Send("Несколько простых советов как быстрее уснуть😴:")
	message := "1. Не ешьте прямо перед тем, как ложиться спать.\n2. Позанимайтесь расслабляющей йогой.\n3. Проветрите комнату." +
		"\n4. Спрячьте часы, чтобы не смотреть на них, пока пытаетесь заснуть.\n5. Уберите телефон подальше от кровати." +
		"\n6. Перед сном примите горячий душ или ванну.\n7. Спите в носках."
	return c.Send(message, menu)
}

func sendAdmin(c tele.Context) error {
	admin, _ := strconv.Atoi(os.Getenv("ADMIN"))
	if c.Sender().ID != int64(admin) {
		return c.Send("Вы не являетесь администратором.")
	}
	userCount, yandexCount := repository.GetAdminInfo()
	return c.Send(fmt.Sprintf("Количество пользователей: %d\nОбщее количество использований Яндекса: %d",
		userCount, yandexCount), keyboards.AdminKB())
}

func massSender(ids []int64, message string, c tele.Context, state fsm.FSMContext) {
	tmp := c.Chat().ID
	for i := 0; i < len(ids); i++ {
		c.Chat().ID = ids[i]
		c.Send(message)
	}
	c.Chat().ID = tmp
}

func yandexToZero(c tele.Context, state fsm.FSMContext) error {
	menu := keyboards.OnStartKB()
	idArray := repository.GetAllId(true)
	massSender(idArray, "Вам обнулилили количество использований Yandex SpeechKit. Можете продолжать использовать его)", c, state)
	return c.Send("Готово, мой господин", menu)
}

func startSendler(cancel tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancel))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(AdminState)
		return c.Send("Введите сообщение для рассылки", menu)

	}
}

func sendler(c tele.Context, state fsm.FSMContext) error {
	message := c.Message().Text
	menu := keyboards.OnStartKB()
	idArray := repository.GetAllId(false)
	massSender(idArray, message, c, state)
	return c.Send("Готово, мой господин", menu)
}
