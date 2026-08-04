# SalesTracker Service

**SalesTracker** — это высокопроизводительный сервисный микросервис на Go для учета, категории и аналитики финансовых транзакций (доходов и расходов) с поддержкой агрегированных метрик, расширенной фильтрации, автоматической выгрузки отчетов в CSV и встроенной веб-панелью управления.

### Технологический стек:

* **Язык**: Go 1.22+
* **База данных**: PostgreSQL (`pgx` / `database/sql`)
* **HTTP-Фреймворк**: Gin / `net/http`
* **Логирование**: Zap / Slog
* **Фронтенд**: HTML5, CSS3, Vanilla JS

---

## Основные возможности

1. **Учет транзакций**: Создание, чтение и удаление записей о доходах (`income`) и расходах (`expense`) с привязкой к категориям и дате.
2. **Фильтрация и Сортировка**: Гибкая выборка операций по типам, категориям и диапазону дат, а также сортировка по ключевым полям.
3. **Финансовая Аналитика**: Расчет метрик в реальном времени:
* Общая сумма и общее количество операций.
* Средний чек (AVG).
* Медиана.
* 90-й перцентиль (P90).
4. **Экспорт данных**: Выгрузка отфильтрованных транзакций в CSV-формате.
5. **Встроенный Dashboard**: Удобный UI для быстрого ввода операций и просмотра аналитических карточек.

---

## Требования к окружению

* **Go**: `1.22` или выше
* **Docker** и **Docker Compose** (для быстрого запуска PostgreSQL)

---

## Запуск и настройка

### 1. Клонирование репозитория

```bash
git clone https://github.com/QuUteO/SalesTracker.git
cd SalesTracker
```

### 2. Запуск инфраструктуры (PostgreSQL)

```bash
docker-compose up -d

```

### 3. Применение миграций БД

Убедитесь, что таблица `items` создана в PostgreSQL:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    type VARCHAR(50) NOT NULL CHECK (type IN ('income', 'expense')),
    category VARCHAR(100) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sales_type ON items(type);
CREATE INDEX IF NOT EXISTS idx_sales_category ON items(category);
CREATE INDEX IF NOT EXISTS idx_sales_date ON items(date);

-- +goose Down
DROP TABLE IF EXISTS items;
```

### 4. Конфигурация (`config.yaml`)

Создайте файл `config.yaml` в корне проекта:

```yaml
env: local

dsn: "postgres://user:user@localhost:5432/salestracker?sslmode=disable"

addr: ":8080"
```

### 5. Запуск сервиса

```bash
go run ./cmd/main.go
```

После запуска сервис доступен по адресу: `http://localhost:8080`

---

## Веб-интерфейс (UI)

Панель управления доступна по адресу: **`http://localhost:8080`**

**Функции UI:**
* Интерактивные карточки метрик (сумма, количество, средний чек, медиана, P90).
* Форма добавления доходов и расходов.
* Динамическая фильтрация таблицы транзакций.
* Кнопка быстрой выгрузки данных в CSV.

---

## API Эндпоинты

### Транзакции (Items / Sales)

#### 1. Добавить транзакцию

* **POST** `/api/v1/items`

* **Body**:

```json
{
  "title": "Оплата сервера",
  "amount": 1500.50,
  "type": "expense",
  "category": "Инфраструктура",
  "date": "2026-08-04T10:00:00Z"
}
```

#### 2. Список транзакций

* **GET** `/api/v1/items`

* **Query параметры**: `type`, `category`, `from`, `to`, `sort_by`, `order`


#### 3. Удалить транзакцию

* **DELETE** `/api/v1/items/:id`


---

### Аналитика и Экспорт

#### 1. Метрики и аналитика

* **GET** `/api/v1/analytics`

* **Ответ (`200 OK`)**:

```json
{
  "count": 42,
  "sum": 125000.00,
  "avg": 2976.19,
  "median": 1800.00,
  "percentile_90": 7500.00
}
```

#### 2. Экспорт в CSV

* **GET** `/api/v1/export/csv`

* Возвращает файл/поток в формате `text/csv`.



---

## Структура проекта

```text
SalesTracker/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── config/              # Загрузка конфигураций
│   ├── handler/             # HTTP Хэндлеры
│   ├── service/             # Бизнес-логика учета и расчета аналитики
│   ├── repository/          # Слой работы с БД (PostgreSQL)
│   └── model/               # DTO и модели данных
├── web/
│   └── index.html           # Dashboard UI
├── docker-compose.yml       # PostgreSQL
├── config.yaml              # Конфиг приложения
├── go.mod
└── README.md
```