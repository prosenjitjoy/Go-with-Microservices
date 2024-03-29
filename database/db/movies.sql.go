// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: movies.sql

package db

import (
	"context"
)

const createMovie = `-- name: CreateMovie :one
INSERT INTO movies (
  id,
  title,
  description,
  director
) VALUES (
  $1, $2, $3, $4
) RETURNING id, title, description, director
`

type CreateMovieParams struct {
	ID          string `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Director    string `db:"director" json:"director"`
}

func (q *Queries) CreateMovie(ctx context.Context, arg *CreateMovieParams) (*Movie, error) {
	row := q.db.QueryRow(ctx, createMovie,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Director,
	)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Director,
	)
	return &i, err
}

const deleteMovie = `-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1
`

func (q *Queries) DeleteMovie(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteMovie, id)
	return err
}

const getMovie = `-- name: GetMovie :one
SELECT id, title, description, director FROM movies
WHERE id = $1
ORDER BY id
LIMIT 1
`

func (q *Queries) GetMovie(ctx context.Context, id string) (*Movie, error) {
	row := q.db.QueryRow(ctx, getMovie, id)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Director,
	)
	return &i, err
}

const listMovies = `-- name: ListMovies :many
SELECT id, title, description, director FROM movies
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListMoviesParams struct {
	Limit  int32 `db:"limit" json:"limit"`
	Offset int32 `db:"offset" json:"offset"`
}

func (q *Queries) ListMovies(ctx context.Context, arg *ListMoviesParams) ([]*Movie, error) {
	rows, err := q.db.Query(ctx, listMovies, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Movie{}
	for rows.Next() {
		var i Movie
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Director,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMovie = `-- name: UpdateMovie :one
UPDATE movies
SET
  title = COALESCE($1, title),
  description = COALESCE($2, description),
  director = COALESCE($3, director)
WHERE
  id = $4
RETURNING id, title, description, director
`

type UpdateMovieParams struct {
	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`
	Director    *string `db:"director" json:"director"`
	ID          string  `db:"id" json:"id"`
}

func (q *Queries) UpdateMovie(ctx context.Context, arg *UpdateMovieParams) (*Movie, error) {
	row := q.db.QueryRow(ctx, updateMovie,
		arg.Title,
		arg.Description,
		arg.Director,
		arg.ID,
	)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Director,
	)
	return &i, err
}
