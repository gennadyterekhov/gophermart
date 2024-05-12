package migration

import (
	"database/sql"
)

func RunMigrations(connection *sql.DB) {
	createTableUsers(connection)
	createTableOrders(connection)
	createTableBalances(connection)
}

func createTableUsers(connection *sql.DB) {
	createTable := `create table if not exists users
	(
		id serial,
		login varchar(255) unique not null,
		password varchar(255) not null
	);`
	_, err := connection.Exec(createTable)
	if err != nil {
		panic(err)
	}
}

func createTableOrders(connection *sql.DB) {
	createTable := `create table if not exists users
	(
		id serial,
		user_id references users(id),
		status varchar(255) not null,
		number varchar(64) not null,
		accrual int default null,
		uploaded_at datetime not null
	);`
	_, err := connection.Exec(createTable)
	if err != nil {
		panic(err)
	}
}

func createTableBalances(connection *sql.DB) {
	createTable := `create table if not exists users
	(
		id serial,
		user_id references users(id),
		current double precision not null,
		withdrawn double precision not null
	);`
	_, err := connection.Exec(createTable)
	if err != nil {
		panic(err)
	}
}
