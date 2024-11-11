package handler

import (
	"context"

	"github.com/justIGreK/MoneyKeeper-Budget/internal/models"
	budgetProto "github.com/justIGreK/MoneyKeeper-Budget/pkg/go/budget"
)

type BudgetServiceServer struct {
	budgetProto.UnimplementedBudgetServiceServer
	BudgetSRV BudgetService
}

type BudgetService interface {
	AddBudget(ctx context.Context, budget models.CreateBudget) (string, error)
	GetBudget(ctx context.Context, userID, budgetID string) (models.Budget, error)
	GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error)
	AddCategory(ctx context.Context, budgetID, catName string, limit float32) (models.Budget, error)
	DeleteCategory(ctx context.Context, budgetID, categoryId string) error
	DeleteBudget(ctx context.Context, budgetID string) error
	UpdateBudget(ctx context.Context, update models.UpdateBudget) (models.Budget, error)
	UpdateCategory(ctx context.Context, update models.) (models.Budget, error)
}

func (s *BudgetServiceServer) AddBudget(ctx context.Context, req *budgetProto.AddBudgetRequest) (*budgetProto.AddBudgetResponse, error) {
	createBudget := models.CreateBudget{
		UserID:    req.UserId,
		Name:      req.Name,
		Limit:     req.Limit,
		Period:    req.Period,
		StartDate: req.Start,
		EndDate:   req.End,
	}
	budgetID, err := s.BudgetSRV.AddBudget(ctx, createBudget)
	if err != nil {
		return nil, err
	}
	return &budgetProto.AddBudgetResponse{
		BudgetId: budgetID,
	}, nil

}

func (s *BudgetServiceServer) GetBudget(ctx context.Context, req *budgetProto.GetBudgetRequest) (*budgetProto.GetBudgetResponse, error) {
	budget, err := s.BudgetSRV.GetBudget(ctx, req.UserId, req.BudgetId)
	if err != nil {
		return nil, err
	}
	protoBudget := budgetProto.Budget{
		BudgetId: budget.ID,
		Name:     budget.Name,
		Limit:    budget.Limit,
		Start:    budget.StartDate.Format(Dateformat),
		End:      budget.EndDate.Format(Dateformat),
		Category: convertToProtoCategories(budget.Category),
	}
	return &budgetProto.GetBudgetResponse{
		Budget: &protoBudget,
	}, nil
}

func convertToProtoCategories(categories []models.Category) []*budgetProto.Category {
	protoBudgets := make([]*budgetProto.Category, len(categories))
	for i, c := range categories {
		protoBudgets[i] = &budgetProto.Category{
			CategoryId: c.ID,
			Name:       c.Name,
			Limit:      c.Limit,
		}
	}
	return protoBudgets
}

func (s *BudgetServiceServer) GetBudgetList(ctx context.Context, req *budgetProto.GetBudgetListRequest) (*budgetProto.GetBudgetListResponse, error) {

	budgets, err := s.BudgetSRV.GetBudgetList(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	protobudgets := convertToProtoBudgets(budgets)
	return &budgetProto.GetBudgetListResponse{
		Budgets: protobudgets,
	}, nil
}

var (
	Dateformat     string = "2006-01-02"
	DateTimeformat string = "2006-01-02T15:04:05"
)

func convertToProtoBudgets(budgets []models.Budget) []*budgetProto.Budget {
	protoBudgets := make([]*budgetProto.Budget, len(budgets))
	for i, b := range budgets {
		protoBudgets[i] = &budgetProto.Budget{
			BudgetId: b.ID,
			Name:     b.Name,
			Limit:    b.Limit,
			Start:    b.StartDate.Format(Dateformat),
			End:      b.EndDate.Format(Dateformat),
			Category: convertToProtoCategories(b.Category),
		}
	}
	return protoBudgets
}
