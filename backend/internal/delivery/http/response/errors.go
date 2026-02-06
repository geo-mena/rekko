package response

import "net/http"

func Error(w http.ResponseWriter, statusCode int, message string) {
	write(w, statusCode, Response{
		Status:  false,
		Message: message,
	})
}

func ValidationError(w http.ResponseWriter, details []ErrorDetail) {
	write(w, http.StatusUnprocessableEntity, Response{
		Status:  false,
		Message: "Error de validación",
		Data: ErrorData{
			Code:    "VALIDATION_ERROR",
			Details: details,
		},
	})
}

func BusinessError(w http.ResponseWriter, statusCode int, message string, code string) {
	write(w, statusCode, Response{
		Status:  false,
		Message: message,
		Data: ErrorData{
			Code: code,
		},
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}
