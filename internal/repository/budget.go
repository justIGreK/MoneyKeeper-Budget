package repository

import (
	"github.com/justIGreK/MoneyKeeper-Budget/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BudgetRepo struct {
	collection *mongo.Collection
}

func NewBudgetRepository(db *mongo.Client) *BudgetRepo {
	return &BudgetRepo{
		collection: db.Database(dbname).Collection(budgetCollection),
	}
}

func (r *BudgetRepo) AddBudget(ctx context.Context, budget models.Budget) (string, error) {
	result, err := r.collection.InsertOne(ctx, budget)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *BudgetRepo) GetBudget(ctx context.Context, userID, budgetID string) (*models.Budget, error) {
	oid, err := convertToObjectIDs(budgetID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var budget models.Budget
	err = r.collection.FindOne(ctx, bson.M{"_id": oid[0], "user_id": userID}).Decode(&budget)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &budget, err
}

func (r *BudgetRepo) GetBudgetList(ctx context.Context, userID string)([]models.Budget, error){
	budgets := []models.Budget{}
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &budgets)
	if err != nil{
		return nil, err
	}
	return budgets, err
}


func (r *BudgetRepo) CloseBudget(ctx context.Context, budgetID string) error {
	oid, err := convertToObjectIDs(budgetID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"is_active": true},
		},
	}

	update := bson.M{"is_active": false}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
