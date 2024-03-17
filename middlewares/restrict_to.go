package middlewares

import (
	"errors"
	"net/http"
	"og-style/models"
	"og-style/utils"
)

type UserRole int

const (
	User UserRole = iota
	Admin
)

func RestrictTo(next http.HandlerFunc, roles ...string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.User)

		for _, role := range user.Role {
			for _, acceptedRole := range roles {
				if role == acceptedRole {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		utils.ForbiddenError(w, errors.New("у вас нет достука к этому маршруту"))
	}
}
