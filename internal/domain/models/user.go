package models

import "time"

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID           int64     `bson:"user_id"`
	Email        string    `bson:"email"`
	PhoneNumber  string    `bson:"phone_number"`
	Name         string    `bson:"name"`
	Surname      string    `bson:"surname"`
	PassHash     string    `bson:"pass_hash"`
	RegisteredAt time.Time `bson:"registered_at"`
	Role         Role      `bson:"role"`
	FamilyIDs    []int64   `bson:"family_ids"`
}
