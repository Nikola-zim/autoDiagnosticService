# autoDiagnosticService
## Общее описание
Система представляет собой HTTP API и телеграмм-бота со следующими
требованиями к бизнес-логике:
- регистрация, аутентификация и авторизация пользователей;
- приём фото и сохранение их в объектное хранилище MinIO;
- отвечать на принимаемые фото в виде изображений с разметкой на
них и текстовой информацией о возможных неисправностях;
- учёт и ведение списка переданных фото и разметки распознанных фото
зарегистрированного пользователя
## Cхема взаимодействия с системой

Первым компонентом является прокси-сервер на GO, который,
используя API телеграмма, получает изображения, создает запросы на их
обработку, получает, обрабатывает и пересылает пользователям результаты.
Второй компонент приложения — это сервер распознавания на Python (https://github.com/Nikola-zim/DetectServerFlask),
на котором происходит распознавание и разметка изображений.
Также компонентом приложения можно назвать телеграмм-бота (<b>@autoDiagnostic_bot</b>), т.к.
он выступает в качестве интерфейса приложения, выполняет функции
аутентификации, идентификации и авторизации пользователей.
Архитектура системы предполагает возможность дублирования
компонентов, для увеличения количества запросов, которые может
обработать система.
Ниже представлена абстрактная бизнес-логика взаимодействия
пользователя с системой:
1. Прокси-сервер на GO получает фото приборной панели от
пользователя. Это может происходить двумя способами: через htmlформу в браузере или через API телеграмма
2. Сохраняем изображение в файловой системе, получив метаданные
объекта, через слой бизнес-логики мы обращаемся к PostgreSQL, где
происходит запись новой строки с полями, которые содержат
метаданные позволяющие в будущем обратиться к объекту в
хранилище. Новый и необработанный запрос имеет пометку «ready» в
поле состояния.
3. Веб-сервис периодически просматривает СУБД в поисках новых
записей с полем «ready». Происходит поиск соответствующих
изображений и передача в модуль распознавания. В начале этого
процесса, отобранным записям присваивается значение «in progress»
4. После обработки этих изображений и получения результатов записям
присваивается значение «recognized»
5. Модуль веб-сервера на Go просматривает записи СУБД с пометкой
«recognized», из этих записей формируются сообщения для ответа и
после их успешной отправки, сервис присваивает соответствующим
записям статус «done».

## Сводное HTTP API
Система должна предоставлять следующие HTTP-хендлеры:
- POST /v1/user/register — регистрация пользователя;
- POST /v1/user/login — аутентификация пользователя;
- POST /v1/user/doRecognition — загрузка пользователем фото для
распознавания;
- GET /v1/user/recognitions — получение истории распознанных
изображений и информации о них;
- GET /v1/user/balance/sum — получение текущего баланса счёта баллов
(кол-во доступных распознаваний);
- POST /v1/user/balance/add — запрос на добавление баллов;
