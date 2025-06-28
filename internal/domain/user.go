package domain

import "time"

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	CreatedDate time.Time `json:"created_date"`
}
