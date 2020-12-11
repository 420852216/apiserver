package user

import "time"

type User struct {
	ID int
	PassWord string
	Name string
	Email string
	Phone string
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}


