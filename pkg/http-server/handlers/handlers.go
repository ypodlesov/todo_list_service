package handlers

import (
	"log/slog"
	"todo_list_service/pkg/storage/postgres"
)

type HandlerContext struct {
	Log     *slog.Logger
	Storage *postgres.Storage
}
