package pg_user

import (
	"awesome-clean-arch/internal/user"
	"awesome-clean-arch/pkg/client/postgresql"
	"awesome-clean-arch/pkg/logging"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
)

type pgRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *pgRepository) Create(ctx context.Context, user user.User) (string, error) {
	q := `INSERT INTO user (username) VALUES ($1) RETURNING id;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if err := r.client.QueryRow(ctx, q, user.Username).Scan(&user.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return "", newErr
		}
		return "", err
	}
	return string(user.ID), nil
}

func (r *pgRepository) FindAll(ctx context.Context) (u []user.User, err error) {
	q := `SELECT id, username FROM awesome.user;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]user.User, 0)

	for rows.Next() {
		var u user.User

		err = rows.Scan(&u.ID, &u.Username)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *pgRepository) FindOne(ctx context.Context, ID string) (user.User, error) {
	q := `SELECT id, username FROM awesome.user WHERE id = $1;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var u user.User

	err := r.client.QueryRow(ctx, q, ID).Scan(u.ID, u.Username)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (r *pgRepository) Update(ctx context.Context, user user.User) error {
	q := `UPDATE awesome.user SET username = $1 WHERE id = $2;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, user.Username, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *pgRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM awesome.user WHERE id = $1;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, ID)
	if err != nil {
		return err
	}

	return nil
}

func NewPGRepository(client postgresql.Client, logger *logging.Logger) user.Repository {
	return &pgRepository{
		client: client,
		logger: logger,
	}
}
