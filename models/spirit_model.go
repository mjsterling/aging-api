package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spirit struct {
	Id         primitive.ObjectID `json:"_id,omitempty"`
	CreatedAt  primitive.DateTime `json:"createdAt,omitempty" validate:"required"`
	Batches    []Batch            `json:"batches,omitempty"`
	Volume     float32            `json:"volume,omitempty" validate:"required"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Type       string             `json:"type,omitempty"`
	InitialABV float32            `json:"initialABV,omitempty" validate:"required"`
	RecipeName string             `json:"recipeName"`
}
