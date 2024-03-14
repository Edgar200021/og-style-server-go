package utils

import (
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, err error, statusCode int) {
	m := map[string]any{
		"status":  "error",
		"message": err.Error(),
	}

	encoded, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(encoded)
}

func SendValidatonErrors(w http.ResponseWriter, errors any) {
	m := map[string]any{
		"status": "error",
		"errors": errors,
	}

	encoded, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(encoded)
}
func BadRequestError(w http.ResponseWriter, err error) {
	SendError(w, err, http.StatusBadRequest)
}
func UnauthorizedError(w http.ResponseWriter, err error) {
	SendError(w, err, http.StatusUnauthorized)
}
func InternalServerError(w http.ResponseWriter, err error) {
	SendError(w, err, http.StatusInternalServerError)
}
func ForbiddenError(w http.ResponseWriter, err error) {
	SendError(w, err, http.StatusForbidden)
}

func SendJSON(w http.ResponseWriter, data any, statusCode int) {
	var m = map[string]any{
		"status": "success",
		"data":   data,
	}

	encoded, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(encoded)
}
