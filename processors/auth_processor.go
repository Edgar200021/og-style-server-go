package processors

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"og-style/db"
	"og-style/types"
	"og-style/utils"
	"time"
)

type AuthProcessor interface {
	SignUp(data types.CreateUser) error
	SignIn(email, password string) (*types.SignInResponse, error)
}

type AuthPgProcessor struct {
	UserStorage  db.UserStorage
	CartStorage  db.CartStorage
	TokenStorage db.TokenStorage
}

func (a *AuthPgProcessor) SignUp(data types.CreateUser) error {
	user, err := a.UserStorage.GetByEmail(data.Email)
	if err != nil {
		return err
	}

	if user.ID != 0 {
		return fmt.Errorf("user with email %s already exists", data.Email)
	}

	hashedPassword, hashErr := utils.HashPassword(data.Password)
	if hashErr != nil {
		return errors.New("something went wrong")
	}

	data.Password = hashedPassword

	userId, createErr := a.UserStorage.Create(&data)
	if createErr != nil {
		return createErr
	}

	if err = a.CartStorage.Create(userId); err != nil {
		return err
	}

	return nil
}
func (a *AuthPgProcessor) SignIn(email, password string) (*types.SignInResponse, error) {

	user, err := a.UserStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("incorrect password or email")
	}

	if ok := utils.CheckPasswordHash(password, user.Password); !ok {
		return nil, errors.New("incorrect password or email")
	}

	var accessToken, refreshToken string

	if dbToken, err := a.TokenStorage.Get(user.ID); err != nil {
		return nil, errors.New("something went wrong")
	} else {
		var res types.SignInResponse
		var tokenError error

		accessToken, tokenError = utils.SignJWT(jwt.MapClaims{
			"id":      user.ID,
			"expires": time.Now().Add(time.Minute * 30),
		})

		if dbToken.ID != 0 {
			refreshToken = dbToken.RefreshToken
		} else {
			refreshToken, _ = utils.SignJWT(jwt.MapClaims{
				"id":      user.ID,
				"expires": time.Now().Add(time.Hour * 24 * 30),
			})
		}

		if tokenError != nil {
			return nil, tokenError
		}

		res.AccessToken = accessToken
		res.RefreshToken = refreshToken
		res.User = user

		return &res, nil
	}
}
