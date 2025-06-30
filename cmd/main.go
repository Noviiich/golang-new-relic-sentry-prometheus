package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/config"
	userHandler "github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/http-server/handlers/user"
	mvLogger "github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/http-server/middleware/logger"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/infrastructure/database"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/lib/logger"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/lib/logger/sl"
	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/repository/postgres"
	userCase "github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)

	log.Info("go new relic sentry prometheus started", slog.String("env", cfg.Env))

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	userRepo := postgres.NewUserRepository(db.DB)
	userUseCase := userCase.NewUserUseCase(userRepo)
	userHandler := userHandler.NewUserHandler(userUseCase, log)

	router := chi.NewRouter()

	// Добавляем middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(mvLogger.New(log))
	// router.Use(response.ErrorHandlerMiddleware(log))
	// router.Use(response.LoggingMiddleware(log))
	// router.Use(response.StatusLoggingMiddleware(log))

	router.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Put("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", userHandler.DeleteUser)
		})
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")

	// // Добавляем health check endpoint
	// router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	response.SendOK(w, r, "Service is healthy", map[string]string{
	// 		"status": "ok",
	// 		"time":   time.Now().Format(time.RFC3339),
	// 	})
	// })
}
