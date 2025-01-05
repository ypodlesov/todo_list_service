package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_list_service/pkg/config"
	"todo_list_service/pkg/http-server/handlers"
	mwLogger "todo_list_service/pkg/http-server/middleware/logger"
	"todo_list_service/pkg/metrics"
	"todo_list_service/pkg/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info(
		"starting todo_list service",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTPServer.Address),
	)

	metrics.StartMetricsServer(&cfg.MetricsConfig)
	storage, err := postgres.New(&cfg.PgConfig)

	if err != nil {
		panic("cannot setup storage")
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	handlerCtx := &handlers.HandlerContext{
		Log:     log,
		Storage: storage,
	}

	router.Post("/sign_up", handlers.NewSignUp(handlerCtx))
	router.Post("/sign_in", handlers.NewSignIn(handlerCtx))
	router.Post("/logout", handlers.NewLogout(handlerCtx))
	router.Get("/get_tasks", handlers.NewGetTasks(handlerCtx))
	router.Get("/get_task", handlers.NewGetTask(handlerCtx))
	router.Post("/create_task", handlers.NewCreateTask(handlerCtx))
	router.Post("/update_task", handlers.NewUpdateTask(handlerCtx))

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
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
		log.Error("failed to stop server", err)
		return
	}

	err = storage.Close()

	if err != nil {
		log.Error("failed to close storage", err)
		return
	}

	log.Info("server stopped")
}
