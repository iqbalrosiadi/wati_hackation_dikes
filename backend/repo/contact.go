package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Contact struct {
	Id    string `json:"id,omitempty" bson:"_id,omitempty"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

type ContactRepo struct {
	col *mongo.Collection
}

func NewContactRepo(col *mongo.Collection) *ContactRepo {
	return &ContactRepo{col}
}

func (r *ContactRepo) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return r.col.Find(ctx, filter)
}

func (r *ContactRepo) FindById(ctx context.Context, id string) (*mongo.SingleResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}
	return r.col.FindOne(ctx, bson.M{"_id": objectId}), nil
}
