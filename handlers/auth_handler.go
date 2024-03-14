package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"og-style/models"
	"og-style/processors"
	"og-style/types"
	"og-style/utils"
	"os"
	"time"
	"unicode/utf8"
)

type AuthHandler struct {
	AuthProcessor processors.AuthPgProcessor
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var data types.CreateUser

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	errors := data.Validate()
	if errors != nil {
		utils.SendValidatonErrors(w, errors)
		return
	}

	if err := a.AuthProcessor.SignUp(data); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "Success", http.StatusCreated)

}
func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var body = make(map[string]string)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.SendError(w, err, http.StatusBadRequest)
		return
	}

	if body["email"] == "" || body["password"] == "" {
		utils.SendError(w, errors.New("all fields are required"), http.StatusBadRequest)
		return
	}

	data, err := a.AuthProcessor.SignIn(body["email"], body["password"])
	if err != nil {
		utils.SendError(w, err, http.StatusBadRequest)
		return
	}

	a.attachTokensToCookie(w, data.AccessToken, data.RefreshToken)
	utils.SendJSON(w, data.User, http.StatusOK)
}
func (a *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		utils.UnauthorizedError(w, errors.New("unauthorized"))
		return
	}

	if data, err := a.AuthProcessor.RefreshTokens(refreshToken.Value); err != nil {
		utils.ForbiddenError(w, err)
		return
	} else {
		a.attachTokensToCookie(w, data.AccessToken, data.RefreshToken)
		utils.SendJSON(w, data.User, http.StatusOK)
	}

}
func (a *AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)
	body := make(map[string]string, 1)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if body["password"] == "" || body["oldPassword"] == "" {
		utils.BadRequestError(w, errors.New("all fields are required"))
		return
	}

	if utf8.RuneCountInString(body["password"]) < 8 {
		utils.BadRequestError(w, errors.New("password must be more or equal 8 characters"))
		return
	}

	if err := a.AuthProcessor.UpdatePassword(user.ID, body["oldPassword"], body["password"]); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "success", http.StatusOK)
}
func (a *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body := make(map[string]string, 1)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if body["email"] == "" {
		utils.BadRequestError(w, errors.New("provide email address"))
		return
	}

	err := a.AuthProcessor.ForgotPassword(body["email"])
	if err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "Check your email for reset password", http.StatusOK)

}
func (a *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	m := make(map[string]string, 1)

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	if m["password"] == "" || utf8.RuneCountInString(m["password"]) < 8 {
		utils.BadRequestError(w, errors.New("password must be more than 8 symbols"))
		return
	}

	if err := a.AuthProcessor.ResetPassword(email, m["password"]); err != nil {
		utils.BadRequestError(w, err)
		return
	}

	utils.SendJSON(w, "success", http.StatusOK)
}
func (a *AuthHandler) attachTokensToCookie(w http.ResponseWriter, accessToken, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		Secure:   os.Getenv("GO_ENV") == "production",
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   os.Getenv("GO_ENV") == "production",
		HttpOnly: true,
	})
}
