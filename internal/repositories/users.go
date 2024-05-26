package repositories

import (
	"context"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"
	"github.com/gennadyterekhov/gophermart/internal/storage"
)

func GetUserById(ctx context.Context, id int64) (*models.User, error) {
	const query = `SELECT id, login, password from users WHERE  id = $1`
	row := storage.DBClient.Connection.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := models.User{}
	err := row.Scan(&(user.ID), &(user.Login), &(user.Password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	const query = `SELECT id, login, password from users WHERE  login = $1`
	row := storage.DBClient.Connection.QueryRowContext(ctx, query, login)
	if row.Err() != nil {
		return nil, row.Err()
	}

	user := models.User{}
	err := row.Scan(&(user.ID), &(user.Login), &(user.Password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func AddUser(ctx context.Context, login string, password string) (*models.User, error) {
	const query = `INSERT INTO users ( login, password) VALUES ( $1, $2) RETURNING id;`

	row := storage.DBClient.Connection.QueryRowContext(ctx, query, login, password)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:       id,
		Login:    login,
		Password: password,
	}

	return &user, nil
}