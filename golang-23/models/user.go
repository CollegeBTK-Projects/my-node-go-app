package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	GroupName string             `bson:"group_name"`
	Course    int                `bson:"course"`
	Grade     int                `bson:"grade"`
	CreatedAt time.Time          `bson:"created_at"`
}
