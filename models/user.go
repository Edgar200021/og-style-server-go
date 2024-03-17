package models

import (
	"time"
)

type User struct {
	ID                   int        `json:"id" db:"id"`
	Email                string     `json:"email" db:"email"`
	Password             string     `json:"-" db:"password"`
	Name                 *string    `json:"name,omitempty" db:"name"`
	Avatar               *string    `json:"avatar,omitempty" db:"avatar"`
	PasswordResetExpires *time.Time `json:"-" db:"password_reset_expires"`
	Role                 []string   `json:"role" db:"role"`
}
