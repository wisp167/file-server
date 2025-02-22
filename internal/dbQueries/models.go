// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package dbQueries

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID         uuid.UUID `json:"id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	FileData   []byte    `json:"file_data"`
	FileName   string    `json:"file_name"`
}
