package response

import (
	"net/http"

	"github.com/go-chi/render"
)

// Response представляет базовую структуру ответа API
type Response struct {
	Code         int         `json:"code,omitempty"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage string      `json:"error,omitempty"`
}

// ErrorResponse представляет структуру для ошибок
type ErrorResponse struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	ErrorMessage string `json:"error,omitempty"`
}

// SuccessResponse представляет структуру для успешных ответов
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse создает новый объект ошибки
func NewErrorResponse(code int, message, errorMessage string) *ErrorResponse {
	return &ErrorResponse{
		Code:         code,
		Message:      message,
		ErrorMessage: errorMessage,
	}
}

// NewSuccessResponse создает новый объект успешного ответа
func NewSuccessResponse(code int, message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// SendError отправляет ошибку клиенту
func SendError(w http.ResponseWriter, r *http.Request, code int, message, errorMessage string) {
	render.Status(r, code)
	render.JSON(w, r, NewErrorResponse(code, message, errorMessage))
}

// SendSuccess отправляет успешный ответ клиенту
func SendSuccess(w http.ResponseWriter, r *http.Request, code int, message string, data interface{}) {
	render.Status(r, code)
	render.JSON(w, r, NewSuccessResponse(code, message, data))
}

// SendBadRequest отправляет ошибку 400
func SendBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	SendError(w, r, http.StatusBadRequest, message, "")
}

// SendInternalServerError отправляет ошибку 500
func SendInternalServerError(w http.ResponseWriter, r *http.Request, message string) {
	SendError(w, r, http.StatusInternalServerError, message, "")
}

// SendNotFound отправляет ошибку 404
func SendNotFound(w http.ResponseWriter, r *http.Request, message string) {
	SendError(w, r, http.StatusNotFound, message, "")
}

// SendCreated отправляет успешный ответ 201
func SendCreated(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	SendSuccess(w, r, http.StatusCreated, message, data)
}

// SendOK отправляет успешный ответ 200
func SendOK(w http.ResponseWriter, r *http.Request, message string, data interface{}) {
	SendSuccess(w, r, http.StatusOK, message, data)
}

func SendNoContent(w http.ResponseWriter, r *http.Request) {
	SendSuccess(w, r, http.StatusNoContent, "", nil)
}
