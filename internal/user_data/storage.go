package user_data

import "context"

type Repository interface {
	Create(ctx context.Context, userData UserData) (string, error)
	FindAll(ctx context.Context) (ud []UserData, err error)
	FindOne(ctx context.Context, userID string) (UserData, error)
	Update(ctx context.Context, userData UserData) error
	Delete(ctx context.Context, userID string) error
}
