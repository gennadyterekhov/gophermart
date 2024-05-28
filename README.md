# gophermart

## how to run

## tests

### coverage
to create coverage file  
`go test -coverprofile=coverage.out ./...`  

to see percentages:
`go tool cover -func=coverage.out`  

to see line by line coverage in browser:  
`go tool cover -html=coverage.out`  


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

## db installation (one-time use)

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
