# ObjectID

Простая библиотека для генерации псевдослучайных глобально уникальных идентификаторов.
Полученные идентификаторы реализуют функцию сравнения, таким образом поддаются сортировке.
Так же из идентификатора можно получить время его генерации, а так же порядковый номер генерации в пределах сессии. Счетчик используемый при генерации идентификаторов потокобезопасный.

## API

```go package objectid // import "go.neonxp.dev/objectid"```


Функции

```go func Seed()```
    необходимо вызвать в начале сессии

```go func New() ID```
    возвращает новый идентификатор

```go func FromString(s string) (ID, error)```
    возвращает идентификатор из base64 представления

```go func FromTime(t time.Time) ID```
    возвращает идентификатор на основе переданного времени

Типы и методы

```go type ID []byte``` тип представляющий собой идентификатор

```go func (i ID) Counter() uint64```
    возвращает порядковый номер идентификатора в сессии

```go func (i ID) Less(i2 ID) bool```
    возвращает true если i2 > i

```go func (i ID) MarshalJSON() ([]byte, error)```
    формирует json представление идентификатора

```go func (i ID) String() string```
    возвращает base64 представление идентификатора

```go func (i ID) Time() time.Time```
    возвращает время создания идентификатора

```go func (i *ID) UnmarshalJSON(b []byte) error```
    парсит идентификатор из json

Примеры

```go
import "go.neonxp.dev/objectid"

objectid.Seed()

id1 := objectid.New()

fmt.Printf("Идентификатор сгенерированный сегодня: %s в %s\n", id1, id1.Time()) // пример: Идентификатор сгенерированный сегодня: AAXwV/DVGwXtTj0FRm92SQF3MiquMPlK в 2022-12-21 18:09:36.872197 +0300 MSK

id2 := objectid.FromTime(time.Now().Add(-24 * time.Hour))

fmt.Printf("Идентификатор сгенерированный вчера: %s в %s\n", id2, id2.Time()) // пример: Идентификатор сгенерированный вчера: AAXwQ+U14N8mbGoVPiiNqyZCss7lEV0Z в 2022-12-20 18:14:42.541791 +0300 MSK

r := "id2 > id1"
if id2.Less(id1) {
    r = "id2 < id1"
}
fmt.Print(r) // выведет: id2 < id1
```

## Лицензия

GNU GPLv3