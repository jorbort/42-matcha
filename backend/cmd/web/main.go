package main

import (
	"log"
	"net/http"
	"context"
	
	"github.com/twpayne/go-geos"
    pgxgeos "github.com/twpayne/pgx-geos"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	serv := http.NewServeMux()

	databaseStr := os.Getenv("DB_URL")
    conn, err := pgx.Connect(context.Background(), databaseStr)
    if err != nil {
        return err
    }
    if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
        return err
    }

	config, err := pgxpool.ParseConfig(connectionStr)
    if err != nil {
        return err
    }
    config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
        if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
            return err
        }
        return nil
    }

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return err
    }
	defer pool.Close()
	
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	serv.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	serv.HandleFunc("GET /{$}", home)

	log.Println("Starting server on :3000")
	err = http.ListenAndServe(":3000", serv)
	log.Fatal(err)
}
