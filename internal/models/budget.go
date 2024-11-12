package models

import "time"

type Budget struct {
	ID        string     `bson:"_id,omitempty"`
	UserID    string     `bson:"user_id"`
	Name      string     `bson:"name"` 
	Limit     float64    `bson:"limit"`
	StartDate time.Time  `bson:"start"`
	EndDate   time.Time  `bson:"end"`
	Category  []Category `bson:"categories"`
}

type Category struct {
	ID    string  `bson:"category_id,omitempty"`
	Name  string  `bson:"name"`
	Limit float64 `bson:"limit"`
}
type CreateBudget struct {
	UserID    string `validate:"required"`
	Name      string `validate:"required"`
	Limit     float64 `validate:"required"`
	Period    string 
	StartDate string
	EndDate   string
}

type CreateCategory struct {
	BudgetID string `validate:"required"`
	UserID    string `validate:"required"`
	Name      string `validate:"required"`
	Limit     float64 `validate:"required"`
}

type GetUpdateBudget struct{
	BudgetID string `validate:"required"`
	UserID string `validate:"required"`
	Name *string 
	Limit *float64 
	Start *string
	End *string
}


type GetUpdateCategory struct{
	BudgetID string `validate:"required"`
	CategoryID string `validate:"required"`
	UserID string `validate:"required"`
	Name *string
	Limit *float64
}

