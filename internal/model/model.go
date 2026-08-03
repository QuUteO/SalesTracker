package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type ItemType string

const (
	TypeIncome  ItemType = "income"
	TypeExpense ItemType = "expense"
)

// Item — основная модель, представляющая запись в БД
type Item struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Amount    float64   `json:"amount" db:"amount"`
	Type      ItemType  `json:"type" db:"type"`
	Category  string    `json:"category" db:"category"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateItemRequest — входящее тело запроса для POST /items
type CreateItemRequest struct {
	Title    string    `json:"title" binding:"required"`
	Amount   float64   `json:"amount" binding:"required,gt=0"`
	Type     ItemType  `json:"type" binding:"required,oneof=income expense"`
	Category string    `json:"category" binding:"required"`
	Date     time.Time `json:"date"`
}

// UpdateItemRequest — входящее тело запроса для PUT /items/{id}
type UpdateItemRequest struct {
	Title    string    `json:"title" binding:"required"`
	Amount   float64   `json:"amount" binding:"required,gt=0"`
	Type     ItemType  `json:"type" binding:"required,oneof=income expense"`
	Category string    `json:"category" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
}

// AnalyticsResponse — ответ для GET /analytics
type AnalyticsResponse struct {
	Count        int64   `json:"count"`
	Sum          float64 `json:"sum"`
	Avg          float64 `json:"avg"`
	Median       float64 `json:"median"`
	Percentile90 float64 `json:"percentile_90"`
}

// AnalyticsFilter — для передачи параметров фильтрации в Service и Repository
type AnalyticsFilter struct {
	From     *time.Time `form:"from"`
	To       *time.Time `form:"to"`
	Category string     `form:"category"`
	Type     string     `form:"type"`
	SortBy   string     `form:"sort_by"` // "date", "amount"
	Order    string     `form:"order"`   // "asc", "desc"
}
