package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/domain"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/lib/api/response"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserCase interface {
	CreateUser(ctx context.Context, user *domain.User) (domain.User, *response.AppError)
	GetUserByID(ctx context.Context, id int) (domain.User, *response.AppError)
	UpdateUser(ctx context.Context, id int, user *domain.User) *response.AppError
	DeleteUser(ctx context.Context, id int) *response.AppError
}

type UserHandler struct {
	uc  UserCase
	log *slog.Logger
}

func NewUserHandler(log *slog.Logger, uc UserCase) *UserHandler {
	return &UserHandler{
		uc:  uc,
		log: log,
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {

		const op = "handlers.user.CreateUser"

		log := h.log.With(
			slog.String("op", op),
			// slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var user domain.User
		if err := render.DecodeJSON(r.Body, &user); err != nil {
			log.Error("failed to decode user", slog.Any("error", err))
			response.SendBadRequest(w, r, "Failed to decode user data")
			return
		}

		user, err := h.uc.CreateUser(r.Context(), &user)
		if err != nil {
			log.Error("failed to create user", slog.Any("error", err))
			hub.CaptureException(errors.New(err.Message))
			response.SendInternalServerError(w, r, "Failed to create user")
			return
		}

		response.SendCreated(w, r, "User created successfully", user)
	}
}

// GetUserByID retrieves a user by ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		const op = "handlers.user.GetUserByID"

		log := h.log.With(
			slog.String("op", op),
			// TODO: add request id
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
			log.Error("invalid user id", slog.String("id", idStr), slog.Any("error", err))
			response.SendBadRequest(w, r, "Invalid user ID format")
			return
		}

		user, getErr := h.uc.GetUserByID(r.Context(), id)
		if getErr != nil {
			log.Error("failed to get user", slog.Int("id", id), slog.Any("error", getErr))
			hub.CaptureException(errors.New(getErr.Message))
			response.SendNotFound(w, r, "User not found")
			return
		}

		response.SendOK(w, r, "User retrieved successfully", user)
	}
}

// UpdateUser updates a user by ID
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		const op = "handlers.user.UpdateUser"

		log := h.log.With(
			slog.String("op", op),
			// TODO: add request id
		)

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Info("user id is required")
			response.SendBadRequest(w, r, "User ID is required")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid user id", slog.String("id", idStr), slog.Any("error", err))
			response.SendBadRequest(w, r, "Invalid user ID format")
			return
		}

		var user domain.User
		if err := render.DecodeJSON(r.Body, &user); err != nil {
			log.Error("failed to decode user", slog.Any("error", err))
			response.SendBadRequest(w, r, "Failed to decode user data")
			return
		}

		if err := h.uc.UpdateUser(r.Context(), id, &user); err != nil {
			log.Error("failed to update user", slog.Int("id", id), slog.Any("error", err))
			hub.CaptureException(errors.New(err.Message))
			response.SendInternalServerError(w, r, "Failed to update user")
			return
		}

		response.SendOK(w, r, "User updated successfully", user)
	}
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		const op = "handlers.user.DeleteUser"

		log := h.log.With(
			slog.String("op", op),
			// TODO: add request id
		)

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Info("user id is required")
			response.SendBadRequest(w, r, "User ID is required")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid user id", slog.String("id", idStr), slog.Any("error", err))
			response.SendBadRequest(w, r, "Invalid user ID format")
			return
		}

		if err := h.uc.DeleteUser(r.Context(), id); err != nil {
			log.Error("failed to delete user", slog.Int("id", id), slog.Any("error", err))
			hub.CaptureException(errors.New(err.Message))
			response.SendInternalServerError(w, r, "Failed to delete user")
			return
		}

		response.SendNoContent(w, r)
	}
}
