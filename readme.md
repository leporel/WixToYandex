# WixToYandex

Конвертация csv выгрузки товаров Wix в формат для шаблона Yandex Market'a  
Создает аналогичный файл csv, от куда можно скопировать товары в шаблон для яндекс маркета.

На коленке написанный за один вечер код, который должен упростить жизнь

## Использование

Посмотрите файл config.toml и настройте под себя  
Запустите программу в консоле с параметром --help

пример запуска

```console
.\wty.exe -f .\catalog_products_2.csv
```

Файл шаблона яндекса - <https://yandex.ru/support/partnermarket/export/excel-format.html>

## Костыли

### Ссылка на товары

Ссылка на товары использует магию космоса, поэтому надеяться что она получиться правильная - не надо Логика такая, что
транслитирируется имя товара и подставляется к той ссылке, что указана в конфигурации, и если небеса будут благосклонны
к вам - ссылка получится рабочей, но если есть дубли, то ссылка точно будет не рабочей, надо проверять ручками.

# По уму, надо сделать (нет)

В wix'е можно зарегать свое приложение и получить доступ к API, авторизация клиентов через OAuth и уже самому делать
выгрузку через апи товаров. Все нужные данные там будут, в отличии от файла csv который предназначен для внутреннего
использования в wix'e