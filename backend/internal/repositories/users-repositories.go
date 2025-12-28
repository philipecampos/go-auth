package repositories

import (
	"auth-jwt/backend/internal/database"
	"auth-jwt/backend/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UsersRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type UsersRepository struct {
	collection *mongo.Collection
}

func NewUsersRepository(db *database.Service) *UsersRepository {
	return &UsersRepository{
		collection: db.Database.Collection("users"),
	}
}

func (r *UsersRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID format: %v", err)
	}

	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objectID}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) Update(ctx context.Context, user *models.User) error {
	if user.ID.IsZero() {
		return errors.New("user ID is required")
	}

	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: user.ID}}, bson.D{{Key: "$set", Value: user}})
	if err != nil {
		return err
	}

	return nil
}
