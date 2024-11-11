package models

import "time"

type Budget struct {
	ID        string     `bson:"_id,omitempty"`
	UserID    string     `bson:"user_id"`
	Name      string     `bson:"name"`
	Limit     float32    `bson:"limit"`
	StartDate time.Time  `bson:"start"`
	EndDate   time.Time  `bson:"end"`
	Category  []Category `bson:"categories,omitempty"`
}

type Category struct {
	ID    string  `bson:"category_id"`
	Name  string  `bson:"name"`
	Limit float32 `bson:"limit"`
}
type CreateBudget struct {
	UserID    string
	Name      string
	Limit     float32
	Period    string
	StartDate string
	EndDate   string
}

type UpdateBudget struct{
	BudgetID string
	UserID string
	Name string
	Limit float32
	Start string
	End string
}

type UpdateCategory struct{
	BudgetID string
	CategoryID string
	UserID string
	Name string
	Limit float32
	Start string
	End string
}
