package response

import (
	"fmt"
	"net/http"
)

// AppError представляет ошибку приложения
type AppError struct {
	Code         int
	Message      string
	ErrorMessage string
	StatusCode   int
}

// Error возвращает строковое представление ошибки
func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, error: %s", e.Code, e.Message, e.ErrorMessage)
}

// NewAppError создает новую ошибку приложения
func NewAppError(code int, message, errorMessage string, statusCode int) *AppError {
	return &AppError{
		Code:         code,
		Message:      message,
		ErrorMessage: errorMessage,
		StatusCode:   statusCode,
	}
}

// CommonErrorCodes определяет общие коды ошибок
const (
	ErrCodeValidationFailed    = 1001
	ErrCodeUserNotFound        = 1002
	ErrCodeUserAlreadyExists   = 1003
	ErrCodeInvalidInput        = 1004
	ErrCodeDatabaseError       = 2001
	ErrCodeInternalServerError = 5001
)

// CommonErrors содержит предопределенные ошибки
var CommonErrors = map[int]*AppError{
	ErrCodeValidationFailed:    NewAppError(ErrCodeValidationFailed, "Validation failed", "", http.StatusBadRequest),
	ErrCodeUserNotFound:        NewAppError(ErrCodeUserNotFound, "User not found", "", http.StatusNotFound),
	ErrCodeUserAlreadyExists:   NewAppError(ErrCodeUserAlreadyExists, "User already exists", "", http.StatusConflict),
	ErrCodeInvalidInput:        NewAppError(ErrCodeInvalidInput, "Invalid input data", "", http.StatusBadRequest),
	ErrCodeDatabaseError:       NewAppError(ErrCodeDatabaseError, "Database operation failed", "", http.StatusInternalServerError),
	ErrCodeInternalServerError: NewAppError(ErrCodeInternalServerError, "Internal server error", "", http.StatusInternalServerError),
}

// GetCommonError возвращает предопределенную ошибку по коду
func GetCommonError(code int) *AppError {
	if err, exists := CommonErrors[code]; exists {
		return err
	}
	return CommonErrors[ErrCodeInternalServerError]
}

// SendAppError отправляет ошибку приложения клиенту
func SendAppError(w http.ResponseWriter, r *http.Request, appErr *AppError) {
	SendError(w, r, appErr.StatusCode, appErr.Message, appErr.ErrorMessage)
}

// SendCommonError отправляет предопределенную ошибку по коду
func SendCommonError(w http.ResponseWriter, r *http.Request, code int) {
	appErr := GetCommonError(code)
	SendAppError(w, r, appErr)
}

func NewNotFoundError(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func NewUnexpectedError(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func NewUserAlreadyExistError(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

func NewValidationError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}
