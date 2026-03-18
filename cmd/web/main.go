package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Arugulamoon/bookkeeper/pkg/config"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
	"github.com/Arugulamoon/bookkeeper/pkg/web"
)

func main() {
	var configFilename string
	flag.StringVar(&configFilename, "config", "", "path to config file")
	flag.Parse()

	if configFilename == "" {
		fmt.Println("missing config filename argument")
		os.Exit(1)
	}

	cfg, err := config.GetConfig(configFilename)
	if err != nil {
		panic(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.DSN())
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templates, err := web.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &web.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,

		DB: db,

		BookAccounts:              &postgres.AccountModel{DB: db},
		BookAssigners:             &postgres.AssignerModel{DB: db},
		BookJournalAccountEntries: &postgres.JournalAccountEntryModel{DB: db},

		SportsRegistrations: &postgres.SportsRegistrationsModel{DB: db},
		SportsMemberships:   &postgres.SportsMembershipModel{DB: db},

		Templates: templates,
	}

	srv := &http.Server{
		Addr:     cfg.Server.Addr(),
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Starting server on %s", cfg.Server.Addr())
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Optional: Configure pool settings (e.g., max connections, lifetime)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse database config: %v", err)
	}
	config.MaxConns = 10
	config.MaxConnLifetime = 30 * time.Minute
	config.MinConns = 2

	// Establish the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}

	return pool, nil
}
