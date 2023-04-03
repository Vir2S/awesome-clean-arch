package mysql_auth

import (
	"awesome-clean-arch/internal/auth"
	"awesome-clean-arch/pkg/client/mysql"
	"awesome-clean-arch/pkg/logging"
	"context"
	"fmt"
	"strconv"
	"strings"
)

type mysqlRepository struct {
	client mysql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *mysqlRepository) Create(ctx context.Context, auth auth.Auth) (string, error) {
	q := `INSERT INTO auth (api_key) VALUES (?);`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	res, err := r.client.ExecContext(ctx, q, auth.APIKey)
	if err != nil {
		r.logger.Error(err)
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		r.logger.Error(err)
		return "", err
	}

	return strconv.FormatInt(id, 10), nil
}

func (r *mysqlRepository) FindAll(ctx context.Context) (u []auth.Auth, err error) {
	q := `SELECT id, api_key FROM auth;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.QueryContext(ctx, q)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	keys := make([]auth.Auth, 0)

	for rows.Next() {
		var a auth.Auth

		err = rows.Scan(&a.ID, &a.APIKey)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		keys = append(keys, a)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return keys, nil
}

func (r *mysqlRepository) FindOne(ctx context.Context, ID string) (auth.Auth, error) {
	q := `SELECT id, api_key FROM auth WHERE id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var a auth.Auth

	err := r.client.QueryRowContext(ctx, q, ID).Scan(&a.ID, &a.APIKey)
	if err != nil {
		r.logger.Error(err)
		return auth.Auth{}, err
	}

	return a, nil
}

func (r *mysqlRepository) Update(ctx context.Context, auth auth.Auth) error {
	q := `UPDATE auth SET api_key = ? WHERE id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, auth.APIKey, auth.ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mysqlRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM auth WHERE id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func NewMySQLRepository(client mysql.Client, logger *logging.Logger) auth.Repository {
	return &mysqlRepository{
		client: client,
		logger: logger,
	}
}
