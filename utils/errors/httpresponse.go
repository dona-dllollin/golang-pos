package errorUtils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	var status int

	switch err {
	case ErrBadRequest:
		status = http.StatusBadRequest
	case ErrConflict:
		status = http.StatusConflict
	case ErrUnauthorized:
		status = http.StatusUnauthorized
	case ErrForbidden:
		status = http.StatusForbidden
	case ErrNotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	log.Printf("[ERROR] %v (status=%d)", err, status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: err.Error(),
	})
}
