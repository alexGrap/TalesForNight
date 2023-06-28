# TalesForNight
Чат-бот сказок на ночь<br/>

## Реализация
Был реализован Telegram Bot, который предоставляет функционал генерации случайной истории по заданому пользователем жанру, или производит генерацию отрывка из указанного произведения.<br/>
Для генерации был использован Chat GPT v3.0, поскольку большая часть мировой литератцры находится под защитой авторского права. Этот метод использован исключительно
в демонстрационных целях, и вместо OpenAi API можно использовать любой иной API ключ с доступом к литературе (требует правок кода). <br/>
Для синтезирования речи использована технология сервиса Yandex SpeechKit. Выбран также в качетсве демонстрации возможностей, и
при желании можно подключить аналог от "Маруси" (powered by VkGroup).<br/>
Для удобство пользователей была добавлена возможность выбора формата предоставляемого текста (Yandex или Текст), выбор жанра для случайно генерируемой истории (Сказка, Поэма или Драма), а также было добавлено меню с интересными фактами про сон и советами по быстроому засыпанию. <br/>
Для удобства получения информации о работе бота был добавлен хендлер /admin, который определенным пользователям предоставляет информации о количестве авторизованных пользователей бота,
возможность добавить всем пользователем количетсво использований Yandex SpeechKit, а также сделать информационную рассылку по всем пользователям.<br/>

---

## Демонстрация
Чтобы "потрогать" бота и ознакомится с его функционалом, перейдите по ссылке: <a href="https://t.me/vkTales_bot">Бот<a/>.<br/>
На данный момент функционирование бота приостановлено. Примеры его использования приведены ниже:
