package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_list_service/internal/config"
	"todo_list_service/internal/http-server/handlers"
	"todo_list_service/internal/http-server/middleware/auth"
	mwLogger "todo_list_service/internal/http-server/middleware/logger"
	"todo_list_service/internal/metrics"
	"todo_list_service/internal/storage/postgres"

	"github.com/gorilla/sessions"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

//	@title			Swagger Example API
//	@version		2.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		todo-list
//	@BasePath	/
func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info(
		"starting todo_list service",
		slog.Any("config", *cfg),
		slog.String("address", cfg.HTTPServer.Address()),
	)

	metrics.StartMetricsServer(&cfg.MetricsConfig)
	storage, err := postgres.New(&cfg.PgConfig)

	logger.Info("created postgres storage")

	if err != nil {
		logger.Error("failed to setup storage", slog.String("error", err.Error()))
		panic("cannot setup storage")
	}

	store := sessions.NewCookieStore([]byte(cfg.Session.SecretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.HTTPServer.Session.MaxAge,
		HttpOnly: true,
		Secure:   cfg.HTTPServer.Session.Secure,
		SameSite: http.SameSiteLaxMode,
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	handlerCtx := &handlers.HandlerContext{
		Log:     logger,
		Storage: storage,
		Store:   store,
	}

	router.Post("/sign_up", handlers.NewSignUp(handlerCtx))
	router.Post("/sign_in", handlers.NewSignIn(handlerCtx))

	authMiddleware := auth.NewAuthMiddleware(store)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Middleware)

		r.Post("/logout", handlers.NewLogout(handlerCtx))
		r.Get("/get_tasks", handlers.NewGetTasks(handlerCtx))
		r.Get("/get_task", handlers.NewGetTask(handlerCtx))
		r.Post("/create_task", handlers.NewCreateTask(handlerCtx))
		r.Post("/update_task", handlers.NewUpdateTask(handlerCtx))
		r.Post("/update_priority", handlers.NewUpdatePriority(handlerCtx))
	})

	logger.Info("starting server", slog.String("address", cfg.HTTPServer.Address()))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("got error, server stopped")
			os.Exit(0)
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", slog.String("error", err.Error()))
		return
	}

	err = storage.Close()

	if err != nil {
		logger.Error("failed to close storage", slog.String("error", err.Error()))
		return
	}

	logger.Info("server stopped")
}
