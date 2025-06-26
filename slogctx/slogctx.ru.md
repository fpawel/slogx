# slogctx

`slogctx` — это пакет для Go, расширяющий стандартный логгер [`log/slog`](https://pkg.go.dev/log/slog) возможностью добавлять атрибуты (ключ-значение) из `context.Context` в записи логов. Это позволяет удобно передавать дополнительные данные через контекст и автоматически включать их в логи.

## Основные возможности

- **Добавление атрибутов в контекст:**  
  Используйте функции `WithValues` и `WithUniqueValues` для добавления пар ключ-значение в контекст.
    - `WithValues` добавляет новые атрибуты, не удаляя дубликаты ключей.
    - `WithUniqueValues` заменяет значения для уже существующих ключей.

- **Удаление атрибутов из контекста:**
    - `WithoutKeys` удаляет указанные ключи из атрибутов контекста.
    - `WithoutAllKeys` полностью очищает все атрибуты.

- **Извлечение атрибутов:**
    - `GetFirstValue` возвращает первое значение по ключу из контекста.
    - `HasKey` проверяет наличие ключа в атрибутах контекста.

- **Интеграция с slog.Handler:**
    - `NewHandler` оборачивает существующий `slog.Handler` и добавляет поддержку атрибутов из контекста.
    - Все атрибуты из контекста автоматически добавляются к каждой записи лога.

## Пример использования

```go
import (
    "context"
    "log/slog"
    "github.com/your/module/slogctx"
)

func main() {
    baseHandler := slog.NewJSONHandler(os.Stdout, nil)
    logger := slog.New(slogctx.NewHandler(baseHandler))

    ctx := slogctx.WithValues(context.Background(), "userID", 42, "requestID", "abc123")
    logger.InfoContext(ctx, "User request received")
}
```

## Документация по функциям

- **NewHandler(handler slog.Handler) slog.Handler**  
  Оборачивает переданный обработчик и добавляет поддержку атрибутов из контекста.

- **WithValues(ctx context.Context, args ...any) context.Context**  
  Добавляет пары ключ-значение в контекст. Ключи должны быть строками.

- **WithUniqueValues(ctx context.Context, args ...any) context.Context**  
  Добавляет пары ключ-значение, заменяя существующие значения для одинаковых ключей.

- **WithoutKeys(ctx context.Context, keys ...string) context.Context**  
  Удаляет указанные ключи из атрибутов контекста.

- **WithoutAllKeys(ctx context.Context) context.Context**  
  Очищает все атрибуты из контекста.

- **GetFirstValue(ctx context.Context, key string) (any, bool)**  
  Возвращает первое значение по ключу и флаг успешного поиска.

- **HasKey(ctx context.Context, key string) bool**  
  Проверяет наличие ключа в атрибутах контекста.

## Особенности и ограничения

- Контекст не предназначен для хранения больших объёмов данных.
- При каждом изменении атрибутов создаётся новый срез, что безопасно для конкурентного доступа, но может быть неэффективно при большом количестве атрибутов.
- Если передано нечётное количество аргументов в `WithValues` или `WithUniqueValues`, последний аргумент игнорируется.
- Дубликаты ключей разрешены в `WithValues` и не фильтруются.

## Лицензия

[MIT](../LICENSE)