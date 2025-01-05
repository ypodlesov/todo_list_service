package main

import (
	"context"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_list_service/pkg/config"
	"todo_list_service/pkg/http-server/handlers"
	"todo_list_service/pkg/http-server/middleware/auth"
	mwLogger "todo_list_service/pkg/http-server/middleware/logger"
	"todo_list_service/pkg/metrics"
	"todo_list_service/pkg/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	log.Info(
		"starting todo_list service",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTPServer.Address),
	)

	metrics.StartMetricsServer(&cfg.MetricsConfig)
	storage, err := postgres.New(&cfg.PgConfig)

	if err != nil {
		log.Error("failed to setup storage", err)
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
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	handlerCtx := &handlers.HandlerContext{
		Log:     log,
		Storage: storage,
		Store:   store,
	}

	router.Post("/sign_up", handlers.NewSignUp(handlerCtx))
	router.Post("/sign_in", handlers.NewSignIn(handlerCtx))
	router.Post("/logout", handlers.NewLogout(handlerCtx))

	authMiddleware := auth.NewAuthMiddleware(store)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Middleware)

		r.Get("/get_tasks", handlers.NewGetTasks(handlerCtx))
		r.Get("/get_task", handlers.NewGetTask(handlerCtx))
		r.Post("/create_task", handlers.NewCreateTask(handlerCtx))
		r.Post("/update_task", handlers.NewUpdateTask(handlerCtx))
	})

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
			log.Error("got error, server stopped")
			os.Exit(0)
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
