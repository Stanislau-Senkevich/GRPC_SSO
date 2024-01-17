package models

import "time"

type User struct {
	Id           int64     `bson:"user_id"`
	Email        string    `bson:"email"`
	PhoneNumber  string    `bson:"phone_number"`
	Name         string    `bson:"name"`
	Surname      string    `bson:"surname"`
	PassHash     string    `bson:"pass_hash"`
	RegisteredAt time.Time `bson:"registered_at"`
}
