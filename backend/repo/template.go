package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Template struct {
	Id       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Language string `json:"language"`
	Content  string `json:"content"`
}

type TemplateRepo struct {
	col *mongo.Collection
}

func NewTemplateRepo(col *mongo.Collection) *TemplateRepo {
	return &TemplateRepo{col}
}

func (r *TemplateRepo) Create(ctx context.Context, template Template) (*mongo.InsertOneResult, error) {
	return r.col.InsertOne(ctx, template)
}

func (r *TemplateRepo) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return r.col.Find(ctx, filter)
}

func (r *TemplateRepo) FindOne(ctx context.Context, id string) *mongo.SingleResult {
	// Convert the string id to MongoDB ObjectID if necessary
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	// Create a filter using BSON
	filter := bson.M{"_id": objectID}

	return r.col.FindOne(ctx, filter)
}
