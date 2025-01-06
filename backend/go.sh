#!/bin/bash

if [ "$DATABASE" = "postgres" ]; then
    echo "Waiting for postgres..."

    while ! nc -z $DATABASE_HOST $SQL_PORT; do
        sleep 0.1
    done

    echo "PostgresSQL started"
fi

#instalo la herramienta goose para las migraciones
go install github.com/pressly/goose/v3/cmd/goose@latest
# aqui se ejecutaran las migraciones de la base de datos cuando las tenga
#goose create add_users_table sql
#goose postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DATABASE_HOST:$SQL_PORT/$POSTGRES_DB?sslmode=disable" up

go run ./cmd/web