package auth

import "context"

type Repository interface {
	Create(ctx context.Context, auth Auth) (string, error)
	FindAll(ctx context.Context) (a []Auth, err error)
	FindOne(ctx context.Context, id string) (Auth, error)
	Update(ctx context.Context, auth Auth) error
	Delete(ctx context.Context, id string) error
}
