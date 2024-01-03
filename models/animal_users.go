package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Collection struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	AnimalId     primitive.ObjectID `json:"animal_id,omitempty" bson:"animal_id,omitempty"`
	DateOfTaking time.Time          `json:"date_of_taking,omitempty" bson:"date_of_taking,omitempty"`
}
