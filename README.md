# gophermart
[ТЗ](https://practicum.yandex.ru/learn/go-advanced/courses/bd3d1b7f-c6b4-4321-bb62-9d246e13ebf1/sprints/261649/topics/e3e4d154-e4ee-4846-9fe4-41f750c16794/lessons/c0d20e78-c1a3-4f49-a13f-733496da984e/)

[DB diagram](https://dbdiagram.io/d/6643a0239e85a46d55d83877)
## how to run

### db installation (one-time use)

      sudo -i -u postgres
      psql -U postgres
      postgres=# create database gophermart_db;
      postgres=# create database gophermart_db_test;
      postgres=# create user gophermart_user with encrypted password 'gophermart_pass';
      postgres=# grant all privileges on database gophermart_db to gophermart_user;
      postgres=# grant all privileges on database gophermart_db_test to gophermart_user;
      alter database gophermart_db owner to gophermart_user;
      alter database gophermart_db_test owner to gophermart_user;
      alter schema public owner to gophermart_user;

after that, use this to connect to db in cli

      psql -U gophermart_user -d gophermart_db

or

      psql -U gophermart_user -d gophermart_db_test

#### migrations

create new migration  

      goose -dir internal/storage/migrations create new_table_users sql
      GOOSE_MIGRATION_DIR="internal/storage/migrations" goose create new_table_orders sql
      GOOSE_MIGRATION_DIR="internal/storage/migrations" goose create new_table_withdrawals sql

run all migrations (real and test db)  

      GOOSE_MIGRATION_DIR="internal/storage/migrations" GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgresql://gophermart_user:gophermart_pass@127.0.0.1:5432/gophermart_db_test?sslmode=disable" goose up



## tests
run 1 suite:  

      go test -run=TestOrders$ ./...

run 1 test method in suite  

      go test -run=TestOrders/Test200IfAlreadyUploaded ./...

run tests and save stdout and stderr to file
go test ./... > .temp/test.log 2>&1


### coverage
to create coverage file  
`go test -coverprofile=coverage.out ./...`  

to see percentages:
`go tool cover -func=coverage.out`  

to see line by line coverage in browser:  
`go tool cover -html=coverage.out`  

### run ci tests locally
      ./ci_tests.sh

# template
# go-musthave-diploma-tpl

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.

## env
JWT_SIGNING_KEY

