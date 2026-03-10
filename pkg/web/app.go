package web

import (
	"database/sql"
	"html/template"
	"log"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	DB *sql.DB

	Templates map[string]*template.Template
}
