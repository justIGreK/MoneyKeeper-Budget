package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/justIGreK/MoneyKeeper-Budget/internal/models"
)

type BudgetRepository interface {
	AddBudget(ctx context.Context, budget models.Budget) (string, error)
	GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error)
	GetBudget(ctx context.Context, userID, budgetID string) (*models.Budget, error)
	AddCategory(ctx context.Context, categ models.CreateCategory) error
	DeleteCategory(ctx context.Context, userID, budgetID, catID string) error
	DeleteBudget(ctx context.Context, userID, budgetID string) error
	UpdateBudget(ctx context.Context, update models.Budget) error
	UpdateCategory(ctx context.Context, userID, budgetID string, update models.Category) error
}

type UserService interface {
	GetUser(ctx context.Context, id string) (string, string, error)
}

type BudgetService struct {
	BudgetRepo BudgetRepository
	User       UserService
}

func NewBudgetService(repo BudgetRepository, user UserService) *BudgetService {
	return &BudgetService{BudgetRepo: repo, User: user}
}

const (
	Dateformat string = "2006-01-02"
)

func (s *BudgetService) AddBudget(ctx context.Context, budget models.CreateBudget) (string, error) {
	user, _, err := s.User.GetUser(ctx, budget.UserID)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if user == "" {
		return "", errors.New("user not found")
	}
	if budget.Limit < 0 {
		budget.Limit *= -1
	}
	var start, end time.Time
	if budget.Period != "" {
		start, end = s.getPeriodDates(budget.Period)
		if start.IsZero() || end.IsZero() {
			return "", errors.New("invalid period")
		}
	} else {
		start, err = time.Parse(Dateformat, budget.StartDate)
		if err != nil {
			log.Println(err)
			return "", err
		}
		end, err = time.Parse(Dateformat, budget.EndDate)
		if err != nil {
			log.Println(err)
			return "", err
		}
		if !start.Before(end) {
			start, end = end, start
		}
	}
	newBudget := models.Budget{
		UserID:    budget.UserID,
		Name:      budget.Name,
		Limit:     float64(budget.Limit),
		StartDate: start,
		EndDate:   end,
		Category:  []models.Category{},
	}
	budgets, err := s.GetBudgetList(ctx, budget.UserID)
	if err != nil {
		return "", err
	}
	for _, budget := range budgets {
		if doTasksOverlap(budget, newBudget) {
			return "", fmt.Errorf("task overlaps with an existing task: %s", budget.Name)
		}
	}
	id, err := s.BudgetRepo.AddBudget(ctx, newBudget)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return id, nil
}

func doTasksOverlap(existingBudget, newBudget models.Budget) bool {
	return existingBudget.EndDate.After(newBudget.StartDate) && existingBudget.StartDate.Before(newBudget.EndDate)
}
func (s *BudgetService) getPeriodDates(period string) (time.Time, time.Time) {
	now := time.Now().UTC()
	switch period {
	case "day":
		start := now
		return start, start.AddDate(0, 0, 1)
	case "week":
		start := now
		return start, start.AddDate(0, 0, 7)
	case "month":
		start := now
		return start, start.AddDate(0, 1, 0)
	case "year":
		start := now
		return start, start.AddDate(1, 0, 0)
	default:
		return time.Time{}, time.Time{}
	}
}
func (s *BudgetService) AddCategory(ctx context.Context, categ models.CreateCategory) (*models.Budget, error) {
	user, _, err := s.User.GetUser(ctx, categ.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	budget, err := s.BudgetRepo.GetBudget(ctx, categ.UserID, categ.BudgetID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if budget == nil {
		return nil, errors.New("budget is not found")
	}
	if s.checkForDuplicateCategory(categ.Name, budget.Category) {
		return nil, fmt.Errorf("category with name %s is already added to this budget", categ.Name)
	}
	if categ.Limit < 0 {
		categ.Limit *= -1
	}

	err = s.BudgetRepo.AddCategory(ctx, categ)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	newBudget, err := s.BudgetRepo.GetBudget(ctx, categ.UserID, categ.BudgetID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return newBudget, nil
}

func (s *BudgetService) checkForDuplicateCategory(newCateg string, categs []models.Category) bool {
	for _, categ := range categs {
		if categ.Name == newCateg {
			return true
		}
	}
	return false
}

func (s *BudgetService) GetBudget(ctx context.Context, userID, budgetID string) (*models.Budget, error) {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	budget, err := s.BudgetRepo.GetBudget(ctx, userID, budgetID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if budget == nil {
		return nil, errors.New("budget is not found")
	}
	return budget, nil
}

func (s *BudgetService) GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error) {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	budgetList, err := s.BudgetRepo.GetBudgetList(ctx, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return budgetList, nil
}

func (s *BudgetService) DeleteCategory(ctx context.Context, userID, budgetID, catID string) error {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if user == "" {
		return errors.New("user not found")
	}
	budget, err := s.BudgetRepo.GetBudget(ctx, userID, budgetID)
	if err != nil {
		log.Println(err)
		return err
	}
	if budget == nil {
		return errors.New("budget is not found")
	}
	isExist := false
	for _, categ := range budget.Category {
		if categ.ID == catID {
			isExist = true
			break
		}
	}
	if !isExist {
		return errors.New("category is not found")
	}

	err = s.BudgetRepo.DeleteCategory(ctx, userID, budgetID, catID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *BudgetService) DeleteBudget(ctx context.Context, userID, budgetID string) error {
	user, _, err := s.User.GetUser(ctx, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	if user == "" {
		return errors.New("user not found")
	}
	budget, err := s.BudgetRepo.GetBudget(ctx, userID, budgetID)
	if err != nil {
		log.Println(err)
		return err
	}
	if budget == nil {
		return errors.New("budget is not found")
	}
	err = s.BudgetRepo.DeleteBudget(ctx, userID, budgetID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *BudgetService) UpdateBudget(ctx context.Context, update models.GetUpdateBudget) (*models.Budget, error) {
	user, _, err := s.User.GetUser(ctx, update.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	budget, err := s.BudgetRepo.GetBudget(ctx, update.UserID, update.BudgetID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if budget == nil {
		return nil, errors.New("budget is not found")
	}
	updates := models.Budget{
		ID:     update.BudgetID,
		UserID: update.UserID,
	}
	if update.Name != nil {
		updates.Name = *update.Name
	} else {
		updates.Name = budget.Name
	}
	if update.Limit != nil {
		updates.Limit = *update.Limit
		if updates.Limit < 0 {
			updates.Limit *= -1
		}
	} else {
		updates.Limit = budget.Limit
	}
	if update.Start != nil {
		updates.StartDate, err = time.Parse(Dateformat, *update.Start)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		updates.StartDate = budget.StartDate
	}
	if update.End != nil {
		updates.EndDate, err = time.Parse(Dateformat, *update.End)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	} else {
		updates.EndDate = budget.EndDate
	}
	if updates.EndDate.Before(updates.StartDate) {
		updates.StartDate, updates.EndDate = updates.EndDate, updates.StartDate
	}
	err = s.BudgetRepo.UpdateBudget(ctx, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	budget, err = s.BudgetRepo.GetBudget(ctx, update.UserID, update.BudgetID)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return budget, nil
}

func (s *BudgetService) UpdateCategory(ctx context.Context, update models.GetUpdateCategory) (*models.Budget, error) {
	user, _, err := s.User.GetUser(ctx, update.UserID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user == "" {
		return nil, errors.New("user not found")
	}
	budget, err := s.BudgetRepo.GetBudget(ctx, update.UserID, update.BudgetID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if budget == nil {
		return nil, errors.New("budget is not found")
	}
	isExist := false
	var existCategory models.Category 
	for _, categ := range budget.Category {
		if categ.ID == update.CategoryID {
			isExist = true
			existCategory = categ
			break
		}
	}
	if !isExist {
		return nil, errors.New("category is not found")
	}
	updates := models.Category{ID: update.CategoryID}
	if update.Name != nil {
		updates.Name = *update.Name
	} else {
		updates.Name = existCategory.Name
	}
	if update.Limit != nil {
		updates.Limit = *update.Limit
	} else {
		updates.Limit = existCategory.Limit
	}
	err = s.BudgetRepo.UpdateCategory(ctx, update.UserID, update.BudgetID, updates)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	budget, err = s.BudgetRepo.GetBudget(ctx, update.UserID, update.BudgetID)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return budget, nil

}
