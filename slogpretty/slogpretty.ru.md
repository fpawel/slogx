# slogpretty

`slogpretty` — это кастомный обработчик для пакета Go `log/slog`, который обеспечивает удобочитаемый, структурированный и цветной вывод логов. Он создан для повышения удобства разработки за счёт более наглядных и информативных логов.

## Возможности

- **Цветной вывод уровней и сообщений** (если вывод в терминал)
- **Настраиваемый формат времени**
- **Поддержка групп атрибутов и переписывания атрибутов**
- **Быстрое форматирование атрибутов** (по умолчанию не использует JSON, но поддерживает его)
- **Опциональный вывод информации о файле и строке**
- **Эффективное использование памяти** через `sync.Pool` для временных срезов

## Использование

### Базовый пример

```go
import (
    "log/slog"
    "os"
    "github.com/fpawel/slogx/slogpretty"
)

func main() {
    handler := slogpretty.NewPrettyHandler().
        WithColorEnabled(true).
        WithTimeLayout("15:04:05").
        WithWriter(os.Stdout)
    logger := slog.New(handler)
    logger.Info("Hello, world!", slog.String("user", "alice"))
}
```

**Пример вывода:**
```
15:04:05 INFO  Hello, world! {"user":"alice"}
```

### Группы атрибутов

```go
logger.Info("Order created",
    slog.Group("order",
        slog.Int("id", 123),
        slog.String("status", "paid"),
    ),
    slog.String("user", "bob"),
)
```
**Пример вывода:**
```
INFO  Order created {"order":{"id":123,"status":"paid"},"user":"bob"}
```

### Кастомный форматтер атрибутов

```go
handler := slogpretty.NewPrettyHandler().
    WithAttrFormatter(func(m map[string]any) string {
        return "ATTRS"
    })
logger := slog.New(handler)
logger.Info("Test", slog.String("foo", "bar"))
```
**Пример вывода:**
```
INFO  Test ATTRS
```

## Конфигурация

- `WithWriter(io.Writer)`: Задать место вывода.
- `WithTimeLayout(string)`: Формат времени.
- `WithLogLevel(slog.Leveler)`: Минимальный уровень логирования.
- `WithSourceInfo(bool)`: Включить/отключить вывод информации о файле и строке.
- `WithAttrRewriter(func)`: Переписать или отфильтровать атрибуты.
- `WithColorEnabled(bool)`: Включить/отключить цветной вывод.
- `WithAttrFormatter(func)`: Кастомное форматирование атрибутов.

## Лицензия

[MIT](../LICENSE)