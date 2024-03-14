package models

type Token struct {
	ID           int    `db:"id"`
	UserID       string `db:"user_id"`
	RefreshToken string `db:"refresh_token"`
}
