package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/justIGreK/MoneyKeeper-Budget/internal/models"

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

func (r *BudgetRepo) GetBudgetList(ctx context.Context, userID string) ([]models.Budget, error) {
	budgets := []models.Budget{}
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &budgets)
	if err != nil {
		return nil, err
	}
	return budgets, err
}

func (r *BudgetRepo) AddCategory(ctx context.Context, categ models.CreateCategory) error {
	oid, err := convertToObjectIDs(categ.BudgetID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	categoryID := primitive.NewObjectID()

	filter := bson.M{"_id": oid[0], "user_id": categ.UserID}
	update := bson.M{"$push": bson.M{"categories": models.Category{
		Name:  categ.Name,
		Limit: categ.Limit,
		ID:    categoryID.Hex(),
	}}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *BudgetRepo) DeleteCategory(ctx context.Context, userID, budgetID, catID string) error {
	oid, err := convertToObjectIDs(budgetID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0], "user_id": userID}
	update := bson.M{"$pull": bson.M{"categories": bson.M{"category_id": catID}}}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *BudgetRepo) DeleteBudget(ctx context.Context, userID, budgetID string) error {
	oid, err := convertToObjectIDs(budgetID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0], "user_id": userID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("budget not found")
	}
	return nil
}

func (r *BudgetRepo) UpdateBudget(ctx context.Context, updates models.Budget) error {
	oid, err := convertToObjectIDs(updates.ID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0], "user_id": updates.UserID}
	update := bson.M{
		"$set": bson.M{
			"name":  updates.Name,
			"limit": updates.Limit,
			"start": updates.StartDate,
			"end":   updates.EndDate,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("UpdateBudget not found")
	}
	return nil
}

func (r *BudgetRepo) UpdateCategory(ctx context.Context, userID, budgetID string, updates models.Category) error {
	oid, err := convertToObjectIDs(budgetID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0], "user_id": userID, "categories.category_id": updates.ID}
	update := bson.M{
		"$set": bson.M{
			"categories.$.name": updates.Name,
			"categories.$.limit":  updates.Limit,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("UpdateCategory error")
	}
	return nil
}
