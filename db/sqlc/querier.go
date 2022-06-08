// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByName(ctx context.Context, username string) (User, error)
}

var _ Querier = (*Queries)(nil)
