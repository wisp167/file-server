-- ./internal/sql

-- name: GetFileByName :many
SELECT * from files
WHERE file_name ILIKE '%' || $1 || '%';

-- name: GetFileByID :one
SELECT * FROM files
WHERE id = $1;

-- name: ListFiles :many
SELECT * FROM files
ORDER BY file_name;

-- name: CreateFile :one
INSERT INTO files (id, file_name, file_data, create_time, update_time)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateFile :one
UPDATE files
SET
    file_name = $2,
    file_data = $3,
    update_time = $4
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;

-- name: CountFiles :one
SELECT COUNT(*) FROM files;

-- name: CountFilesByName :one
SELECT COUNT(*) FROM files
WHERE file_name ILIKE '%' || $1 || '%';

-- name: SearchFiles :many
SELECT * FROM files
WHERE file_name ILIKE '%' || $1 || '%';
