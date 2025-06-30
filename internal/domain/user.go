package domain

import "time"

type User struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Age         int       `json:"age" db:"age"`
	CreatedDate time.Time `json:"created_date" db:"created_date"`
}
