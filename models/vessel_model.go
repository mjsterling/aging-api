package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vessel struct {
	Id        primitive.ObjectID `json:"_id,omitempty"`
	Batches   []Batch            `json:"batches,omitempty"`
	CreatedAt primitive.DateTime `json:"createdAt,omitempty" validate:"required"`
	Volume    float32            `json:"volume,omitempty" validate:"required"`
	Material  string             `json:"material,omitempty" validate:"material"`
	Process   string             `json:"process" validate:"process"`
}
