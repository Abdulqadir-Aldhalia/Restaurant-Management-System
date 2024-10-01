package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("conflict occurred")
	ErrInvalidArgument = errors.New("invalid request")
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	log.Printf("Error: %v", err)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		log.Printf("Error: Type = ErrNoRows \n Details = %s\n", err.Error())
		statusCode = http.StatusNotFound
		message = "Resource not found"

	case errors.Is(err, ErrNotFound):
		log.Printf("Error: Type = ErrNotFound \n Details = %s\n", err.Error())
		statusCode = http.StatusNotFound
		message = "Requested resource not found"

	case errors.Is(err, ErrConflict):
		log.Printf("Error: Type = ErrConflict \n Details = %s\n", err.Error())
		statusCode = http.StatusConflict
		message = "Conflict occurred"

	case errors.Is(err, ErrInvalidArgument):
		log.Printf("Error: Type = ErrInvalidArgument \n Details = %s\n", err.Error())
		statusCode = http.StatusBadRequest
		message = "Invalid request"

	default:
		log.Printf("Error: Type = InternalServerError \n Details = %s\n", err.Error())
		statusCode = http.StatusInternalServerError
		message = "Oops, something went wrong"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}

func SendCustomeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}
