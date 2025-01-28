package storage

import (
	"context"
	"github.com/wisp167/file-server/internal/dbQueries"
	"github.com/google/uuid"
)

type Storage interface {
	Select(ctx context.Context, columns []string) ([]dbQueries.File, error)
	GetFile(ctx context.Context, id uuid.UUID) (dbQueries.File, error)
}
