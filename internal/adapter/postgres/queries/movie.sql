-- name: UpsertMovie :exec
INSERT INTO movies (id, title, year, director)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET title    = EXCLUDED.title,
    year     = EXCLUDED.year,
    director = EXCLUDED.director;

-- name: GetMovie :one
SELECT id, title, year, director
FROM movies
WHERE id = $1;

-- name: GetMoviesRange :many
SELECT id, title, year, director
FROM movies
WHERE sqlc.arg(from_id) <= id AND id <= sqlc.arg(to_id)
ORDER BY id;

-- name: GetLastMovie :one
SELECT id, title, year, director
FROM movies
ORDER BY id DESC
LIMIT 1;

-- name: ExistIDs :many
SELECT id FROM movies;
