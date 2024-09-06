package repo

import "go.mongodb.org/mongo-driver/mongo"

type CompiledContactProfile struct {
	Id        string   `json:"id" bson:"_id"`
	Phone     string   `json:"phone"`
	Contact   string   `json:"contact"`
	Location  string   `json:"location"`
	Age       int32    `json:"age"`
	Interests []string `json:"interests"`
}

type CompiledContactProfileRepo struct {
	col *mongo.Collection
}

func NewCompiledContactProfileRepo(col *mongo.Collection) *CompiledContactProfileRepo {
	return &CompiledContactProfileRepo{col}
}
