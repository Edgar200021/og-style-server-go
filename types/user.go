package types

type CreateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type UpdateUser struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=8"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}
