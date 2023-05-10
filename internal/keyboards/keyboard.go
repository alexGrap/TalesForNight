package keyboards

import tele "gopkg.in/telebot.v3"

var (
	FindBtn   = tele.Btn{Text: "🔍 Найти сказку"}
	SpeechBtn = tele.Btn{Text: "🎙 Задать формат"}
	GenreBtn  = tele.Btn{Text: "🖋 Выбрать жанр"}
	InfoBtn   = tele.Btn{Text: "🕶 Дополнительная информация"}
	UserBtn   = tele.Btn{Text: "⚙ Ваши настройки"}

	OwnTaleBtn = tele.Btn{Text: "📔 Выбрать свою книгу"}
	OurTaleBtn = tele.Btn{Text: "📚 Выбрать случайную книгу"}

	YandexBtn = tele.Btn{Text: "✨ Yandex SpeechKit"}

	FairyBtn = tele.Btn{Text: "🎆 Сказка"}
	PoemBtn  = tele.Btn{Text: "✒ Поэма"}
	DramaBtn = tele.Btn{Text: "🎭 Драма"}

	GetSleepingInfoBtn = tele.Btn{Text: "📝 Интересная инфомрация о сне"}
	SleepingAdviceBtn  = tele.Btn{Text: "✌ Советы"}

	TextBtn = tele.Btn{Text: "🧾 Текст"}

	CancelBtn = tele.Btn{Text: "❌ Назад"}

	AdminSendlerBtn = tele.Btn{Text: "Сделать рассылку сообщений"}
	AdminYandexBtn  = tele.Btn{Text: "Обнулить использование Яндекса"}
)

func OnStartKB() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(FindBtn, SpeechBtn),
		menu.Row(GenreBtn, InfoBtn), menu.Row(UserBtn))
	return menu
}

func AdminKB() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(AdminSendlerBtn, AdminYandexBtn),
		menu.Row(CancelBtn))
	return menu
}
