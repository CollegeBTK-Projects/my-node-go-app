package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"    json:"id"`
	Name      string        `bson:"name"             json:"name"       binding:"required,min=2,max=100"`
	Email     string        `bson:"email"            json:"email"      binding:"required,email"`
	GroupName string        `bson:"group_name"       json:"group_name" binding:"required,group"`
	Age       int           `bson:"age,omitempty"    json:"age"        binding:"omitempty,gte=16,lte=100"`
	Course    int           `bson:"course,omitempty" json:"course,omitempty" binding:"omitempty,gte=1,lte=4"`
	CreatedAt time.Time     `bson:"createdAt"        json:"createdAt"`
}
