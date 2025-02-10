package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jorbort/42-matcha/user_managment/internals/models"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
)

type aplication struct {
	models    *models.Models
	templates *template.Template
}

func main() {

	ctx := context.Background()

	pool, err := createDb(os.Getenv("DB_URL"), ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer pool.Close()

	templates := initTempaltes()

	app := &aplication{
		models:    &models.Models{DB: pool},
		templates: templates,
	}

	if err := app.models.CreateTables(ctx); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Starting server on :3000")
	err = http.ListenAndServe(":3000", app.routes())
	log.Fatal(err.Error())
}

func createDb(dns string, ctx context.Context) (*pgxpool.Pool, error) {

	conn, err := pgx.Connect(context.Background(), dns)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
		log.Fatal(err.Error())
	}

	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		log.Fatal(err.Error())
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
			return err
		}
		return nil
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return pool, nil
}

func initTempaltes() *template.Template {
	templates := template.Must(template.ParseFiles("ui/html/header_template.html",
		"ui/html/profile.html",
		"ui/html/validation.html"))
	return templates
}
