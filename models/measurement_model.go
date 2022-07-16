package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Measurement struct {
	ID         primitive.ObjectID `json:"_id,omitempty"`
	CreatedAt  primitive.DateTime `json:"createdAt,omitempty" validate:"required"`
	ABV        float32            `json:"abv,omitempty" validate:"required"`
	Image      string             `json:"image,omitempty" validate:"required"`
	Nose       string             `json:"nose,omitempty"`
	ForePalate string             `json:"forePalate,omitempty"`
	MidPalate  string             `json:"midPalate,omitempty"`
	Finish     string             `json:"finish,omitempty"`
	Notes      string             `json:"notes,omitempty"`
}