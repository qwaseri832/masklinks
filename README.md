# masklinks

[![CI](https://github.com/qwaseri832/masklinks/actions/workflows/ci.yml/badge.svg)](https://github.com/qwaseri832/masklinks/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)

Библиотека и CLI, которые скрывают адреса ссылок в тексте, оставляя видимой
только схему. Полезно, когда текст надо показать или залогировать, но сами
адреса светить нельзя.

```
Hello, its my page: http://localhost123.com See you
Hello, its my page: http://**************** See you
```

## Установка

```bash
go install github.com/qwaseri832/masklinks/cmd/masklinks@latest
```

Либо из исходников:

```bash
git clone https://github.com/qwaseri832/masklinks.git
cd masklinks
go build ./cmd/masklinks
```

## Использование

Текст берётся из аргументов, а если их нет — из stdin:

```bash
masklinks "мой сайт: http://example.com"
```

```bash
cat notes.txt | masklinks > notes.masked.txt
```

По умолчанию обрабатываются схемы `http://` и `https://`. Список можно задать
своим:

```bash
masklinks -schemes "ftp://,sftp://" "архив на ftp://files.local/dump.tar"
```

## Как библиотека

```go
import "github.com/qwaseri832/masklinks"

masked := masklinks.Mask("см. http://example.com/a?b=1 дальше")
// "см. http://***************** дальше"

masked = masklinks.MaskSchemes("ftp://files.local", "ftp://")
// "ftp://***********"

masked = masklinks.MaskSchemes(text, append(masklinks.DefaultSchemes(), "ftp://")...)
```

`DefaultSchemes` возвращает копию, поэтому дополнять её `append` безопасно.

## Правила маскирования

| Правило | Пример |
|---|---|
| Схема сохраняется как есть, включая регистр | `HTTP://Example.COM` → `HTTP://***********` |
| Адресом считается всё до первого пробельного символа | пробел, `\t`, `\n` завершают адрес |
| Считаются символы, а не байты | `http://сайт.рф` → `http://*******` (7 звёздочек, не 13) |
| Схема без адреса не трогается | `http:// пусто` → `http:// пусто` |
| Длина текста в символах не меняется | инвариант проверяется тестом |

Сохранение длины — причина, по которой маскирование написано вручную, а не
регуляркой с заменой на фиксированную строку: разметка текста не съезжает,
а сам адрес по результату не восстановить.

## Разработка

```bash
go test -race ./...
go vet ./...
gofmt -l .
```

## Лицензия

[MIT](LICENSE)
