package middlewares

import (
	"context"
	"errors"
	"net/http"
	"og-style/db"
	"og-style/utils"
)

func Auth(handler http.HandlerFunc, userStorage db.UserStorage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, tokenErr := r.Cookie("accessToken")
		if tokenErr != nil {
			utils.UnauthorizedError(w, errors.New("unauthorized"))
			return
		}

		claims, err := utils.ParseJWT(accessToken.Value)
		if err != nil {
			utils.UnauthorizedError(w, err)
			return
		}

		if user, err := userStorage.Get(int(claims["id"].(float64))); err != nil {
			utils.UnauthorizedError(w, err)
			return
		} else {
			if user.ID == 0 {
				utils.UnauthorizedError(w, err)
				return
			}
			ctx := context.WithValue(context.Background(), "user", user)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
