package processors

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"og-style/db"
	"og-style/types"
	"og-style/utils"
	"time"
)

type AuthProcessor interface {
	SignUp(data types.CreateUser) error
	SignIn(email, password string) (*types.SignInResponse, error)
	RefreshTokens(refreshToken string) (*types.SignInResponse, error)
	UpdatePassword(userId int, oldPassword, password string) error
	ForgotPassword(email string) error
	ResetPassword(email, password string) error
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
		return fmt.Errorf("пользователь с эл.почтой %s уже существует", data.Email)
	}

	hashedPassword, hashErr := utils.HashPassword(data.Password)
	if hashErr != nil {
		return errors.New("что-то пошло не так.Повторите попытку чуть позже")
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
		return nil, errors.New("неправильный пароль или эл.адрес")
	}

	if ok := utils.CheckPasswordHash(password, user.Password); !ok {
		return nil, errors.New("неправильный пароль или эл.адрес")
	}

	var accessToken, refreshToken string

	if dbToken, err := a.TokenStorage.Get(user.ID); err != nil {
		return nil, errors.New("что-то пошло не так.Повторите попытку чуть позже")
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

			if err := a.TokenStorage.Create(user.ID, refreshToken); err != nil {
				return nil, err
			}
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
func (a *AuthPgProcessor) RefreshTokens(refreshToken string) (*types.SignInResponse, error) {

	token, err := utils.ParseJWT(refreshToken)
	if err != nil {
		return nil, err
	}

	if user, err := a.UserStorage.Get(int(token["id"].(float64))); err != nil {
		return nil, err
	} else {
		if user.ID == 0 {
			return nil, err
		}

		dbToken, err := a.TokenStorage.Get(user.ID)
		if err != nil {
			return nil, err
		}

		if dbToken.RefreshToken != refreshToken {
			return nil, err
		}

		if accessToken, err := utils.SignJWT(jwt.MapClaims{
			"id":      user.ID,
			"expires": time.Now().Add(time.Minute * 30),
		}); err != nil {
			return nil, err
		} else {
			return &types.SignInResponse{
				User:         user,
				AccessToken:  accessToken,
				RefreshToken: dbToken.RefreshToken,
			}, nil
		}
	}
}
func (a *AuthPgProcessor) UpdatePassword(userId int, oldPassword, password string) error {

	user, err := a.UserStorage.Get(userId)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return fmt.Errorf("user with id %d doesn't exists", userId)
	}

	if ok := utils.CheckPasswordHash(oldPassword, user.Password); !ok {
		return errors.New("invalid password")
	}

	if hashedPassword, err := utils.HashPassword(password); err != nil {
		return err
	} else {
		if err := a.UserStorage.UpdatePassword(userId, hashedPassword); err != nil {
			return err
		}
		return nil
	}

}
func (a *AuthPgProcessor) ForgotPassword(email string) error {
	emailCh := make(chan error)
	updateUserCh := make(chan error)

	user, err := a.UserStorage.GetByEmail(email)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return fmt.Errorf("пользователь с эл.почтой %s не существует", email)
	}

	templateParser, tempErr := template.ParseFiles("./templates/forgot-password.html")
	if tempErr != nil {
		return tempErr
	}

	buff := &bytes.Buffer{}
	if err := templateParser.Execute(buff, struct{ Email, Token string }{
		Email: email,
	}); err != nil {
		return err
	}

	go func() {
		if err := utils.SendEmail(buff.String(), "Reset password", email); err != nil {
			emailCh <- err
			return
		}
		emailCh <- nil
	}()

	go func() {
		if err := a.UserStorage.UpdatePasswordExpires(user.ID, time.Now().Add(time.Minute*15)); err != nil {
			updateUserCh <- err
			return
		}
		updateUserCh <- nil
	}()

	if emailErr, updateUserErr := <-emailCh, <-updateUserCh; emailErr != nil || updateUserErr != nil {
		fmt.Println(emailErr)
		fmt.Println(updateUserErr)
		return errors.New("что-то пошло не так.Повторите попытку чуть позже")
	}
	return nil
}
func (a *AuthPgProcessor) ResetPassword(email, password string) error {
	deleteExpiresCh := make(chan error)
	updatePasswordCh := make(chan error)

	user, err := a.UserStorage.GetByEmail(email)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return fmt.Errorf("пользователь с эл.почтой %s не существует", email)
	}

	fmt.Println(user)
	if time.Now().After(user.PasswordResetExpires.Add(time.Second * 0)) {
		return errors.New("время восстановления пароля истек")
	}

	if hashedPassword, err := utils.HashPassword(password); err != nil {
		return err
	} else {
		go func() {
			if err := a.UserStorage.UpdatePassword(user.ID, hashedPassword); err != nil {
				updatePasswordCh <- err
				return
			}
			updatePasswordCh <- nil
		}()
		go func() {
			if err := a.UserStorage.DeletePasswordResetExpires(user.ID); err != nil {
				deleteExpiresCh <- err
				return
			}
			deleteExpiresCh <- nil
		}()

		if deleteErr, updateErr := <-deleteExpiresCh, <-updatePasswordCh; deleteErr != nil || updateErr != nil {
			fmt.Println(deleteErr)
			fmt.Println(updateErr)
			return errors.New("что-то пошло не так.Повторите попытку чуть позже")
		}
		return nil
	}

}
