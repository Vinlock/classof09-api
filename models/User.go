package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Name string `bson:"name,omitempty" json:"name,omitempty"`

	Email string `bson:"email,omitempty" json:"email,omitempty"`

	FacebookId string `bson:"facebookId,omitempty" json:"facebookId,omitempty"`
}
