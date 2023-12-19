## Задание
1. Вытащить из апи Центробанка (пример http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=11/11/2020) данные по переводу различных валют в рубли за последние 90 дней.
2. Результатом работы программы:  
 - нужно указать значение максимального курса валюты, название этой валюты и дату этого максимального значения.
 - нужно указать значение минимального курса валюты, название этой валюты и дату этого минимального значения.
 - нужно указать среднее значение курса рубля за весь период по всем валютам.

### Установка
```
git clone https://github.com/rawenfalls/goCbrf.git
```
### Сборка
```
make
```
### Запуск Linux
```
./apiWork
```
### Запуск Windows
```
apiWork.exe
```
