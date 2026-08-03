-- +goose Up
CREATE TYPE type_sales AS ENUM('доходы', 'расходы');

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    amount NUMERIC NOT NULL CHECK (amount >= 0),
    type type_sales,
    category VARCHAR(100),
    date TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP NOT NULL  DEFAULT now()
);

-- 1. Индекс для фильтрации и сортировки по дате
CREATE INDEX idx_items_date ON items(date DESC);

-- 2. Индекс для фильтрации по категории
CREATE INDEX idx_items_category ON items(category);

-- 3. Составной индекс для быстрой аналитики по типам операций за период
CREATE INDEX idx_items_type_date ON items(type, date);


-- +goose Down
DROP TABLE IF EXISTS items;
