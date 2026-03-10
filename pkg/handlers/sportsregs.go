package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type SportsRegistrationHandler struct {
	DB *sql.DB
}

type ListData struct {
	Type string `query:"type"`
}

func (h *SportsRegistrationHandler) ListRegistrations() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(ListData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		sportsRegModel := &postgres.SportsRegistrationsModel{DB: h.DB}
		var regs []*models.SportsRegistration
		var err error
		switch data.Type {
		case "currentAndUpcoming":
			regs, err = sportsRegModel.SelectAllCurrentAndUpcoming()
		case "past":
			regs, err = sportsRegModel.SelectAllPast()
		default:
			regs, err = sportsRegModel.SelectAll()
		}
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "sportsregs.list.page.tmpl",
			map[string]any{
				"Type":                data.Type,
				"SportsRegistrations": regs,
			})
	}
}
