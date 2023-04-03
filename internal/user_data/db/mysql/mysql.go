package mysql_user_data

import (
	"awesome-clean-arch/internal/user_data"
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

func (r *mysqlRepository) Create(ctx context.Context, ud user_data.UserData) (string, error) {
	q := `INSERT INTO user_data (school) VALUES (?);`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	res, err := r.client.ExecContext(ctx, q, ud.School)
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

func (r *mysqlRepository) FindAll(ctx context.Context) (ud []user_data.UserData, err error) {
	q := `SELECT user_id, school FROM user_data;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.QueryContext(ctx, q)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	data := make([]user_data.UserData, 0)

	for rows.Next() {
		var d user_data.UserData

		err = rows.Scan(&d.ID, &d.School)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		data = append(data, d)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (r *mysqlRepository) FindOne(ctx context.Context, ID string) (user_data.UserData, error) {
	q := `SELECT user_id, school FROM user_data WHERE user_id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var ud user_data.UserData

	err := r.client.QueryRowContext(ctx, q, ID).Scan(&ud.ID, &ud.School)
	if err != nil {
		r.logger.Error(err)
		return user_data.UserData{}, err
	}

	return ud, nil
}

func (r *mysqlRepository) Update(ctx context.Context, ud user_data.UserData) error {
	q := `UPDATE user_data SET school = ? WHERE user_id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, ud.School, ud.ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mysqlRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM user_data WHERE user_id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func NewMySQLRepository(client mysql.Client, logger *logging.Logger) user_data.Repository {
	return &mysqlRepository{
		client: client,
		logger: logger,
	}
}
