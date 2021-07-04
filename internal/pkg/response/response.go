package response

import (
	"encoding/json"
	"fmt"
	"github.com/mineway/logger"
	"net/http"
	"os"
	"reflect"
)

type content struct {
	Message string `json:"message"`
}

var errorMask = "an error has occurred"

// Error allows to print JSON error/text response
func Error(w http.ResponseWriter, status int, err interface{}) bool {
	var errorText, errorValue string

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if reflect.ValueOf(err).Type().String() == "string" {
		errorText = err.(string)
	} else {
		errorText = err.(error).Error()
	}

	if status == http.StatusInternalServerError {
		if os.Getenv("DISPLAY_ERROR") != "true" {
			errorValue = errorMask
		} else {
			errorValue = errorText
		}
		logger.Warning("[API][InternalServerError] %s", errorText)
	}else{
		errorValue = errorText
	}

	failed := json.NewEncoder(w).Encode(&content{ Message: errorValue })
	if failed != nil {
		return false
	}

	return true
}

// Success allows to print any type of structure in json return
func Success (w http.ResponseWriter, status int, success interface{}) bool{
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if isNil(success) {
		_, _ = fmt.Fprintf(w, "{}")
		return true
	}

	failed := json.NewEncoder(w).Encode(success)
	if failed != nil {
		return false
	}

	return true
}

// SuccessText allows to print string in json return
func SuccessText (w http.ResponseWriter, status int, message string) bool{
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	failed := json.NewEncoder(w).Encode(&content{ Message: message })
	if failed != nil {
		return false
	}

	return true
}

// NoContent send 204 response
func NoContent (w http.ResponseWriter) bool {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	return true
}

// Detect if an array is empty
func isNil(i interface{}) bool {
	if reflect.ValueOf(i).Kind() == reflect.Slice {
		return reflect.ValueOf(i).IsNil()
	}
	return false
}