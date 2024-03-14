package models

type Cart struct {
	ID     int `db:"id"`
	UserID int `db:"user_id"`
}
