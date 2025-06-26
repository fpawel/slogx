# slogtest

`slogtest` — это пакет Go, предоставляющий утилиты для тестирования кода, использующего стандартный пакет `log/slog`. Он помогает перехватывать, проверять и утверждать вывод логов в модульных тестах.

## Возможности

- **Обработчик логов в памяти**: сохраняет все записи логов для последующей проверки.
- **Удобное создание тестовых логгеров**: быстрое получение логгера для тестов.
- **Вывод логов без временных меток**: для детерминированного вывода в тестах.

## Использование

### Базовый пример

```go
import (
    "testing"
    "log/slog"
    "github.com/fpawel/slogx/slogtest"
)

func TestLogging(t *testing.T) {
    logger, observed := slogtest.NewTestLogger(t)
    logger.Info("hello", slog.String("key", "value"))
    logs := observed.Logs()
    if len(logs) != 1 || logs[0].Message != "hello" {
        t.Errorf("unexpected logs: %+v", logs)
    }
}
```

### Вывод логов без временных меток

```go
import (
    "log/slog"
    "github.com/fpawel/slogx/slogtest"
)

func main() {
    handler := slogtest.NewStdoutTextHandlerWithoutTimestamp()
    logger := slog.New(handler)
    logger.Info("no timestamp here")
}
```

## API

- `NewTestLogger(t *testing.T) (*slog.Logger, *ObservedHandler)`: возвращает логгер и обработчик для проверки логов в тестах.
- `NewStdoutTextHandlerWithoutTimestamp() slog.Handler`: возвращает обработчик, выводящий логи в stdout без временных меток.
- `ObservedHandler.Logs() []ObservedLog`: возвращает все перехваченные логи для анализа.

## Лицензия

[MIT](../LICENSE)