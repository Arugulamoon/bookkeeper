package web

import (
	"html/template"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger

	DB *pgxpool.Pool

	BookAccounts              *postgres.AccountModel
	BookAssigners             *postgres.AssignerModel
	BookJournalAccountEntries *postgres.JournalAccountEntryModel

	SportsRegistrations *postgres.SportsRegistrationsModel
	SportsMemberships   *postgres.SportsMembershipModel

	Templates map[string]*template.Template
}
