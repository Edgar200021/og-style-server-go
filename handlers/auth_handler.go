package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"og-style/processors"
	"og-style/types"
	"og-style/utils"
	"os"
	"time"
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
