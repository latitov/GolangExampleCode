# GolangExampleCode
Example/reference code, in Go language

"Покажи мне код!". Примеры кода, задача которых уменьшать количество сказанных слов и экономить время.

Данные файлы являются __примерами продакшн кода__, вырванными из контекста реальных проектов, разных проектов. Каждый файл призван продемонстрировать тот или иной юзкейс.

АННОТАЦИЯ ПО ФАЙЛАМ:

- HTTP_WS_Server.go
    
    Самый "аутистический код", без декомпозиции. Почему? Все спрашивали, здесь отвечаю:
    1. Код __сгенерирован__, фактический исходник имеет декомпозицию на куски по 20-40 строк, на уровне специализированного текстового редактора.
    1. Зоны ответственности по исключениям __разделены__, каждый хендлер обрабатывает свою сам, если падает весь сервер то сигнализирует в главный коммутатор (main interconnect) наверх с помощью сваливания контекста.
    1. Magic numbers подставлены автоматически на этапе сборки (билда).
    1. Юнит-тестами да не покрыто, увы, честно. Если бы правильно, то так бы не прокатило, кучей. Time shortage to blame.

- DataModel_S1.go
    
    Построение модели данных для внутреннего меж-горутинного взаимодействия. Как пример работы со сложными иерархическими структурами. Декомпозиция (отсутствующая) реализована на уровне текстового редактора.

- __MODBUSProtocolProcessor.go__

    Уже декомпозиция на месте, и даже юнит тесты есть (здесь не приложены).
    
    Чем интересен данный код, это тем что здесь мы идим работу со встроенным брокером сообщений. "Встроенный", это масштаба процесса, между горутинами, произвольным их количеством.

- jdmqll.go

    "Нормальный" код. За шаблон взята container/list.go из старндартной библиотеки, и модифицировано для использования в качестве JDMQLL. Что это, то отдельная история.

- __SQL_DB_nav.go__

    Тоже нормальный код, не аутистический. С декомпозицией. Однако, комменты вычищены и вырвано из контекста. Тем не менее, что это? Это библиотека функций для работы с БД PostgreSQL, реализация специализированного иерархического дерева под специализированную учетную систему, поверх реляционной модели. В конце есть универсальный инструмент выполнения произвольных SQL программ с созданием транзакции. Просто пример работы с реляционной БД.

Мне посоветовали лайф-хак, взять проекты с наибольшим количеством звезд, и сделать что-то подобное - как образец. Согласен. Но пока не сделал, пока не было времени. Пока это. Будет время - сделаю. Пока так.
