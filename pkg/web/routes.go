package web

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/Arugulamoon/bookkeeper/pkg/handlers"
	"github.com/Arugulamoon/bookkeeper/pkg/sports"
)

func (app *Application) Routes() http.Handler {
	e := echo.New()

	// HTML Template Renderer
	e.Renderer = &TemplateRenderer{
		templates: app.Templates,
	}

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Secure())

	// Routes
	annualReportHandler := &handlers.AnnualReportHandler{
		BookAccounts:              app.BookAccounts,
		BookJournalAccountEntries: app.BookJournalAccountEntries,
	}
	e.GET("/", annualReportHandler.Show())

	e.GET("/annualreport", annualReportHandler.Show())

	jAcctHandler := &handlers.JournalAccountHandler{
		BookAccounts:              app.BookAccounts,
		BookJournalAccountEntries: app.BookJournalAccountEntries,
	}
	e.GET("/journal/account/:acctType/:acctName", jAcctHandler.Show())

	jAcctEntryHandler := &handlers.JournalAccountEntryHandler{
		BookAccounts:              app.BookAccounts,
		BookJournalAccountEntries: app.BookJournalAccountEntries,
	}
	e.GET("/journal/account/entries/:id/edit", jAcctEntryHandler.EditForm())
	e.POST("/journal/account/entries/:id/edit", jAcctEntryHandler.Edit())
	e.GET("/journal/account/entries/:id/split", jAcctEntryHandler.SplitForm())
	e.POST("/journal/account/entries/:id/split", jAcctEntryHandler.Split())

	assignerHandler := &handlers.AssignerHandler{
		BookAccounts:              app.BookAccounts,
		BookAssigners:             app.BookAssigners,
		BookJournalAccountEntries: app.BookJournalAccountEntries,
	}
	e.GET("/assigners", assignerHandler.List())
	e.GET("/assigner/create", assignerHandler.CreateForm())
	e.POST("/assigner/create", assignerHandler.Create())
	e.GET("/assigner/:id", assignerHandler.Show())

	accountHandler := &handlers.AccountHandler{
		BookAccounts:              app.BookAccounts,
		BookJournalAccountEntries: app.BookJournalAccountEntries,
	}
	e.GET("/accounts", accountHandler.List())
	e.GET("/account/create", accountHandler.CreateForm())
	e.POST("/account/create", accountHandler.Create())
	e.GET("/account/:acctType/:acctName", accountHandler.Show())

	sportsRegHandler := &handlers.SportsRegistrationHandler{
		SportsRegistrations: app.SportsRegistrations,
	}
	e.GET("/sports/registrations", sportsRegHandler.ListRegistrations())

	sportsMembershipHandler := &handlers.SportsMembershipHandler{
		SportsMemberships: app.SportsMemberships,
	}
	e.GET("/sports/memberships", sportsMembershipHandler.List())
	e.GET("/sports/memberships/:id", sportsMembershipHandler.Show())

	sports2RegistrationsHandler := &sports.RegistrationsHandler{DB: app.DB}
	e.GET("/sports2/registrations", sports2RegistrationsHandler.ListRegistrations())

	return e
}
