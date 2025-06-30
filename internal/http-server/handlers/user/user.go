package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/domain"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/lib/api/response"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserCase interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	DeleteUser(id int) error
}

type UserHandler struct {
	uc  UserCase
	log *slog.Logger
}

func NewUserHandler(uc UserCase, log *slog.Logger) *UserHandler {
	return &UserHandler{
		uc:  uc,
		log: log,
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// if hub := sentry.GetHubFromContext(r.Context()); hub != nil {

	const op = "handlers.user.CreateUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var user domain.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode user", sl.Err(err))
		response.SendBadRequest(w, r, "Failed to decode user data")
		return
	}

	user, err := h.uc.CreateUser(user)
	if err != nil {
		log.Error("failed to create user", sl.Err(err))
		// hub.CaptureException(errors.New(err.Error()))
		response.SendInternalServerError(w, r, "Failed to create user")
		return
	}

	response.SendCreated(w, r, "User created successfully", user)
	// 	}
}

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
	const op = "handlers.user.GetUserByID"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	// Получаем ID из URL параметров
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Info("user id is required")
		response.SendBadRequest(w, r, "User ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid user id", slog.String("id", idStr), sl.Err(err))
		response.SendBadRequest(w, r, "Invalid user ID format")
		return
	}

	user, getErr := h.uc.GetUserByID(id)
	if getErr != nil {
		log.Error("failed to get user", slog.Int("id", id), sl.Err(getErr))
		// hub.CaptureException(errors.New(getErr.Error()))
		response.SendNotFound(w, r, "User not found")
		return
	}

	response.SendOK(w, r, "User retrieved successfully", user)
	// }
}

// UpdateUser updates a user by ID
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
	const op = "handlers.user.UpdateUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var user domain.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode user", sl.Err(err))
		response.SendBadRequest(w, r, "Failed to decode user data")
		return
	}

	user, err := h.uc.UpdateUser(user)
	if err != nil {
		log.Error("failed to update user", sl.Err(err))
		// hub.CaptureException(errors.New(err.Error()))
		response.SendInternalServerError(w, r, "Failed to update user")
		return
	}

	response.SendOK(w, r, "User updated successfully", user)
	// }
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
	const op = "handlers.user.DeleteUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Info("user id is required")
		response.SendBadRequest(w, r, "User ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid user id", slog.String("id", idStr), sl.Err(err))
		response.SendBadRequest(w, r, "Invalid user ID format")
		return
	}

	if err := h.uc.DeleteUser(id); err != nil {
		log.Error("failed to delete user", slog.Int("id", id), sl.Err(err))
		// hub.CaptureException(errors.New(err.Error()))
		response.SendInternalServerError(w, r, "Failed to delete user")
		return
	}

	response.SendNoContent(w, r)
	// }
}
