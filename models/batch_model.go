package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Batch struct {
	Id           primitive.ObjectID `json:"_id,omitempty"`
	CreatedAt    primitive.DateTime `json:"createdAt,omitempty" validate:"required"`
	Vessels      []Vessel           `json:"vessels,omitempty"`
	Measurements []Measurement      `json:"measurements,omitempty"`
	Volume       float32            `json:"volume,omitempty" validate:"required"`
}
