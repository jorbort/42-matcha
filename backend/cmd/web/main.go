package main

import (
	"log"
	"net/http"
	"context"
	"os"
	
	"github.com/jorbort/42-matcha/backend/internals/models"
	"github.com/twpayne/go-geos"
    pgxgeos "github.com/twpayne/pgx-geos"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type aplication struct {
	models *models.Models
}

func main() {
	
	ctx := context.Background()
	
    pool, err := createDb(os.Getenv("DB_URL"), ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer pool.Close()
	app := &aplication{models: &models.Models{DB: pool}}
	
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