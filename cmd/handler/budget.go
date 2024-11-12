package handler

import (
	"context"
	"errors"

	"github.com/go-playground/validator"
	"github.com/justIGreK/MoneyKeeper-Budget/internal/models"
	budgetProto "github.com/justIGreK/MoneyKeeper-Budget/pkg/go/budget"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BudgetServiceServer struct {
	budgetProto.UnimplementedBudgetServiceServer
	BudgetSRV BudgetService
}

type BudgetService interface {
	AddBudget(ctx context.Context, budget models.CreateBudget) (string, error)
	GetBudget(ctx context.Context, userID, budgetID string) (*models.Budget, error)
	GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error)
	AddCategory(ctx context.Context, categ models.CreateCategory) (*models.Budget, error)
	DeleteCategory(ctx context.Context, userID, budgetID, categoryId string) error
	DeleteBudget(ctx context.Context, userID, budgetID string) error
	UpdateBudget(ctx context.Context, update models.GetUpdateBudget) (*models.Budget, error)
	UpdateCategory(ctx context.Context, update models.GetUpdateCategory) (*models.Budget, error)
}

var validate = validator.New()

func (s *BudgetServiceServer) AddBudget(ctx context.Context, req *budgetProto.AddBudgetRequest) (*budgetProto.AddBudgetResponse, error) {
	createBudget := models.CreateBudget{
		UserID:    req.UserId,
		Name:      req.Name,
		Limit:     float64(req.Limit),
		Period:    req.Period,
		StartDate: req.Start,
		EndDate:   req.End,
	}
	if err := validate.Struct(createBudget); err != nil {
		return nil, err
	}
	if createBudget.Period == "" && (createBudget.StartDate == "" && createBudget.EndDate == "") {
		return nil, errors.New("missing agruments: period")
	}
	budgetID, err := s.BudgetSRV.AddBudget(ctx, createBudget)
	if err != nil {
		return nil, err
	}
	return &budgetProto.AddBudgetResponse{
		BudgetId: budgetID,
	}, nil

}

func (s *BudgetServiceServer) AddCategory(ctx context.Context, req *budgetProto.AddCategoryRequest) (*budgetProto.GetBudgetResponse, error) {
	addCategory := models.CreateCategory{
		UserID:   req.UserId,
		BudgetID: req.BudgetId,
		Name:     req.Name,
		Limit:    float64(req.Limit),
	}
	if err := validate.Struct(addCategory); err != nil {
		return nil, err
	}
	budget, err := s.BudgetSRV.AddCategory(ctx, addCategory)
	if err != nil {
		return nil, err
	}
	protoBudget := budgetProto.Budget{
		BudgetId: budget.ID,
		Name:     budget.Name,
		Limit:    float32(budget.Limit),
		Start:    budget.StartDate.Format(Dateformat),
		End:      budget.EndDate.Format(Dateformat),
		Category: convertToProtoCategories(budget.Category),
	}
	return &budgetProto.GetBudgetResponse{
		Budget: &protoBudget,
	}, nil

}

func (s *BudgetServiceServer) DeleteCategory(ctx context.Context, req *budgetProto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	err := s.BudgetSRV.DeleteCategory(ctx, req.UserId, req.BudgetId, req.CategoryId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BudgetServiceServer) DeleteBudget(ctx context.Context, req *budgetProto.DeleteBudgetRequest) (*emptypb.Empty, error) {
	err := s.BudgetSRV.DeleteBudget(ctx, req.UserId, req.BudgetId)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BudgetServiceServer) UpdateBudget(ctx context.Context, req *budgetProto.UpdateBudgetRequest) (*budgetProto.GetBudgetResponse, error) {
	updateBudget := models.GetUpdateBudget{
		BudgetID: req.Update.BudgetId,
		UserID:   req.Update.UserId,
	}
	if err := validate.Struct(updateBudget); err != nil {
		return nil, err
	}
	if err := s.validateUpdateBudget(req); err != nil{
		return nil, err 
	}
	if req.Update.Name != nil {
		updateBudget.Name = &req.Update.Name.Value 
	}
	if req.Update.Limit != nil {
		updateBudget.Limit = &req.Update.Limit.Value 
	}
	if req.Update.Start != nil {
		updateBudget.Start = &req.Update.Start.Value 
	}
	if req.Update.End != nil {
		updateBudget.End = &req.Update.End.Value 
	}
	
	budget, err := s.BudgetSRV.UpdateBudget(ctx, updateBudget)
	if err != nil {
		return nil, err
	}
	protoBudget := budgetProto.Budget{
		BudgetId: budget.ID,
		Name:     budget.Name,
		Limit:    float32(budget.Limit),
		Start:    budget.StartDate.Format(Dateformat),
		End:      budget.EndDate.Format(Dateformat),
		Category: convertToProtoCategories(budget.Category),
	}
	return &budgetProto.GetBudgetResponse{
		Budget: &protoBudget,
	}, nil
}
func (s *BudgetServiceServer)validateUpdateBudget(req *budgetProto.UpdateBudgetRequest) error {
	if req.Update.Name == nil && req.Update.Limit == nil &&
	req.Update.Start == nil && req.Update.End == nil{
		return errors.New("no new updates")
	}
	return nil
}

func (s *BudgetServiceServer)validateUpdateCategory(req *budgetProto.UpdateCategoryRequest) error {
	if req.Update.Name == nil && req.Update.Limit == nil {
		return errors.New("either 'Name' or 'Limit' must be provided")
	}
	return nil
}


func (s *BudgetServiceServer) UpdateCategory(ctx context.Context, req *budgetProto.UpdateCategoryRequest) (*budgetProto.GetBudgetResponse, error) {
	updateCategory :=  models.GetUpdateCategory{
		BudgetID:   req.Update.BudgetId,
		UserID:     req.Update.UserId,
		CategoryID: req.Update.CategoryId,
	}
	if err := validate.Struct(updateCategory); err != nil {
		return nil, err
	}
	if err := s.validateUpdateCategory(req); err != nil{
		return nil, err
	}
	if req.Update.Name != nil {
		updateCategory.Name = &req.Update.Name.Value 
	}
	if req.Update.Limit != nil {
		updateCategory.Limit = &req.Update.Limit.Value 
	}
	budget, err := s.BudgetSRV.UpdateCategory(ctx, updateCategory)
	if err != nil {
		return nil, err
	}
	
	protoBudget := budgetProto.Budget{
		BudgetId: budget.ID,
		Name:     budget.Name,
		Limit:    float32(budget.Limit),
		Start:    budget.StartDate.Format(Dateformat),
		End:      budget.EndDate.Format(Dateformat),
		Category: convertToProtoCategories(budget.Category),
	}
	return &budgetProto.GetBudgetResponse{
		Budget: &protoBudget,
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
		Limit:    float32(budget.Limit),
		Start:    budget.StartDate.Format(Dateformat),
		End:      budget.EndDate.Format(Dateformat),
		Category: convertToProtoCategories(budget.Category),
	}
	return &budgetProto.GetBudgetResponse{
		Budget: &protoBudget,
	}, nil
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
			Limit:    float32(b.Limit),
			Start:    b.StartDate.Format(Dateformat),
			End:      b.EndDate.Format(Dateformat),
			Category: convertToProtoCategories(b.Category),
		}
	}
	return protoBudgets
}

func convertToProtoCategories(categories []models.Category) []*budgetProto.Category {
	protoBudgets := make([]*budgetProto.Category, len(categories))
	for i, c := range categories {
		protoBudgets[i] = &budgetProto.Category{
			CategoryId: c.ID,
			Name:       c.Name,
			Limit:      float32(c.Limit),
		}
	}
	return protoBudgets
}
