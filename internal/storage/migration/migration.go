package migration

import (
	"database/sql"
)

func RunMigrations(connection *sql.DB) {
	createTableUsers(connection)
	createTableOrders(connection)
	createTableWithdrawals(connection)
}

func createTableUsers(connection *sql.DB) {
	createTable := `create table if not exists users
	(
		id serial not null primary key,
		login varchar(255) unique not null,
		password varchar(255) not null
	);`
	_, err := connection.Exec(createTable)
	if err != nil {
		panic(err)
	}
}

func createTableOrders(connection *sql.DB) {
	createTable := `create table if not exists orders
	(
		number varchar(64) unique not null primary key,
		user_id int references users(id),
		status varchar(255) not null,
		accrual int default null,
		uploaded_at timestamp with time zone not null default now()
	);`
	_, err := connection.Exec(createTable)
	if err != nil {
		panic(err)
	}
}

func createTableWithdrawals(connection *sql.DB) {
	createTable := `create table if not exists withdrawals
	(
		id serial not null primary key,
		user_id int references users(id),
		order_number int not null,
    	total_sum int not null,
		processed_at timestamp with time zone not null
	);`
	_, err := connection.Exec(createTable)
	if err != nil {
		panic(err)
	}
}
