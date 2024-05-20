package repositories

import (
	"context"
	"database/sql"

	"github.com/gennadyterekhov/gophermart/internal/logger"

	"github.com/gennadyterekhov/gophermart/internal/domain/models"

	"github.com/gennadyterekhov/gophermart/internal/storage"
)

func GetUserById(ctx context.Context, id int) (*models.User, error) {
	tx, err := storage.Connection.DBConnection.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		err := tx.Commit()
		if err != nil {
			logger.ZapSugarLogger.Errorln(err.Error())
		}
	}(tx)

	if err != nil {
		return nil, err
	}
	return GetUserByIdTx(ctx, tx, id)
}

func GetUserByIdTx(ctx context.Context, tx *sql.Tx, id int) (*models.User, error) {
	const query = `select id, login, password from users where id = $1`
	row := tx.QueryRowContext(ctx, query, id)
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

// AddUserTx expects already encrypted password
func AddUserTx(ctx context.Context, tx *sql.Tx, login string, password string) (*models.User, error) {
	const query = `insert into users ( login, password) values ( $1, $2) RETURNING id;`

	row := tx.QueryRowContext(ctx, query, login, password)
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
