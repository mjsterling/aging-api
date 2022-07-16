package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Batch struct {
	ID           primitive.ObjectID `json:"_id,omitempty"`
	Vessels      []Vessel           `json:"vessels,omitempty"`
	Measurements []Measurement      `json:"measurements,omitempty"`
	Volume       float32            `json:"volume,omitempty" validate:"required"`
}
