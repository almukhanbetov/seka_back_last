package models

type Driver struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Email  string `db:"email" json:"email"`
	Image  string `db:"image" json:"image"`
	Status int    `db:"status" json:"status"`
}
