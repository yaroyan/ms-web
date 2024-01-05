package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	repository "github.com/yaroyan/ms/authn/infrastructure/persistence/repository/postgres"
	"github.com/yaroyan/ms/authn/interfaces/api/handler"
	"github.com/yaroyan/ms/authn/interfaces/api/router"
)

const port = 80

func main() {
	log.Printf("Starting authentication service on port %d.\n", port)

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can not connect to Postgres.")
	}

	r := router.Router{
		AuthnHandler: &handler.AuthenticationHandler{
			Repository: &repository.UserRepository{
				Connection: conn,
			},
			Client: &http.Client{},
		},
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r.Routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres.")
			return connection
		}

		if 10 < counts {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds ...")
		time.Sleep(2 * time.Second)
		continue
	}
}
