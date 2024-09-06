package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Broadcast struct {
	Id         string    `json:"id" bson:"_id"`
	Name       string    `json:"name"`
	TemplateId string    `json:"template_id"`
	Content    string    `json:"content"`
	Contacts   []Contact `json:"contacts"`
}

type BroadcastRepo struct {
	col *mongo.Collection
}

func NewBroadcastRepo(col *mongo.Collection) *BroadcastRepo {
	return &BroadcastRepo{col}
}

func (r *BroadcastRepo) Create(ctx context.Context, broadcast Broadcast) error {
	_, err := r.col.InsertOne(ctx, broadcast)
	return err
}

func (r *BroadcastRepo) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return r.col.Find(ctx, filter)
}
