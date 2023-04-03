package profile

import "context"

type Repository interface {
	Create(ctx context.Context, profile Profile) (string, error)
	FindAll(ctx context.Context) (p []Profile, err error)
	FindOne(ctx context.Context, userID string) (Profile, error)
	Update(ctx context.Context, profile Profile) error
	Delete(ctx context.Context, userID string) error
}
