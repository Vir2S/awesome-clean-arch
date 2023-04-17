package mysql_profile

import (
	"awesome-clean-arch/internal/profile"
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

func (r *mysqlRepository) Create(ctx context.Context, p profile.Profile) (string, error) {
	q := `START TRANSACTION;
INSERT INTO user (username) VALUES (?);
SET @user_id = LAST_INSERT_ID();
INSERT INTO user_profile (user_id, first_name, last_name, phone, address, city) 
VALUES (@user_id, ?, ?, ?, ?, ?);
INSERT INTO user_data (user_id, school) VALUES (@user_id, ?);
COMMIT;
`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	r.logger.Infoln("p.ID = ", p.ID, "p.Username = ", p.Username, "p.FirstName = ", p.FirstName,
		"p.LastName = ", p.LastName, "p.Phone = ", p.Phone, "p.Address = ", p.Address,
		"p.City = ", p.City, "p.School = ", p.School)

	res, err := r.client.ExecContext(ctx, q, p.Username, p.FirstName, p.LastName, p.Phone, p.Address, p.City, p.School)
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

func (r *mysqlRepository) FindAll(ctx context.Context) (p []profile.Profile, err error) {
	q := `SELECT user.username, user_profile.user_id, user_profile.first_name,
       user_profile.last_name, user_profile.phone, user_profile.address, user_profile.city, user_data.school
	FROM user JOIN user_profile ON user.id = user_profile.user_id JOIN user_data ON user.id = user_data.user_id;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.QueryContext(ctx, q)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	profiles := make([]profile.Profile, 0)

	for rows.Next() {
		var up profile.Profile

		err = rows.Scan(&up.Username, &up.ID, &up.FirstName, &up.LastName, &up.Phone, &up.Address, &up.City, &up.School)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}

		profiles = append(profiles, up)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err)
		return nil, err
	}

	return profiles, nil
}

func (r *mysqlRepository) FindOne(ctx context.Context, Username string) (profile.Profile, error) {
	q := `SELECT user.username, user_profile.user_id, user_profile.first_name, user_profile.last_name,
       user_profile.phone, user_profile.address, user_profile.city, user_data.school
	FROM user JOIN user_profile ON user.id = user_profile.user_id JOIN user_data ON user.id = user_data.user_id WHERE user.username = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var up profile.Profile

	err := r.client.QueryRowContext(ctx, q, Username).Scan(&up.Username, &up.ID, &up.FirstName, &up.LastName, &up.Phone, &up.Address, &up.City, &up.School)
	if err != nil {
		r.logger.Error(err)
		return profile.Profile{}, err
	}

	return up, nil
}

func (r *mysqlRepository) Update(ctx context.Context, p profile.Profile) error {
	q := `UPDATE user_profile SET first_name = ?, last_name = ?, phone = ?, address = ?, city = ? WHERE user_id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, p.FirstName, p.LastName, p.Phone, p.Address, p.City, p.ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mysqlRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM user_profile WHERE user_id = ?;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.ExecContext(ctx, q, ID)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func NewMySQLRepository(client mysql.Client, logger *logging.Logger) profile.Repository {
	return &mysqlRepository{
		client: client,
		logger: logger,
	}
}
