package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vessel struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Batches   []Batch            `json:"batches"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	Volume    float32            `json:"volume,omitempty" validate:"required"`
	Material  string             `json:"material,omitempty" validate:"material"`
	Process   string             `json:"process" validate:"process"`
}
