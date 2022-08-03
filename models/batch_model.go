package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Batch struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt    primitive.DateTime `json:"createdAt"`
	Vessels      []Vessel           `json:"vessels"`
	Measurements []Measurement      `json:"measurements"`
	Volume       float32            `json:"volume,omitempty" validate:"required"`
}
