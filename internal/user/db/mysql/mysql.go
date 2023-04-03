package mysql_user

import (
	"awesome-clean-arch/internal/user"
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

func (r *mysqlRepository) Create(ctx context.Context, user user.User) (string, error) {
	q := `INSERT INTO user (username) VALUES (?);`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	res, err := r.client.ExecContext(ctx, q, user.Username)
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

func (r *mysqlRepository) FindAll(ctx context.Context) (u []user.User, err error) {
	q := `SELECT id, username FROM user;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.QueryContext(ctx, q)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	users := make([]user.User, 0)

	for rows.Next() {
		var u user.User

		err = rows.Scan(&u.ID, &u.Username)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return users, nil
}

func (r *mysqlRepository) FindOne(ctx context.Context, ID string) (user.User, error) {
	q := `SELECT id, username FROM user WHERE id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var u user.User

	err := r.client.QueryRowContext(ctx, q, ID).Scan(&u.ID, &u.Username)
	if err != nil {
		r.logger.Error(err)
		return user.User{}, err
	}

	return u, nil
}

func (r *mysqlRepository) Update(ctx context.Context, user user.User) error {
	q := `UPDATE user SET username = ? WHERE id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, user.Username, user.ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mysqlRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM user WHERE id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func NewMySQLRepository(client mysql.Client, logger *logging.Logger) user.Repository {
	return &mysqlRepository{
		client: client,
		logger: logger,
	}
}
