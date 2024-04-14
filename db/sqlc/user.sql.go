// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(
    username, password, full_name, email, avatar
)VALUES (
    $1, $2, $3, $4, $5
) RETURNING username, password, full_name, email, avatar, password_changed_at, created_at
`

type CreateUserParams struct {
	Username string         `json:"username"`
	Password string         `json:"password"`
	FullName string         `json:"full_name"`
	Email    string         `json:"email"`
	Avatar   sql.NullString `json:"avatar"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Password,
		arg.FullName,
		arg.Email,
		arg.Avatar,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Password,
		&i.FullName,
		&i.Email,
		&i.Avatar,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, password, full_name, email, avatar, password_changed_at, created_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.Password,
		&i.FullName,
		&i.Email,
		&i.Avatar,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
