package models

import "time"

type User struct {
	ID        int
	Name      string
	GroupName string
	Course    string
	Grade     float32
	CreatedAt time.Time
}
