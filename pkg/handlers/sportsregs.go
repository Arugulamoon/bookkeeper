package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type SportsRegistrationHandler struct {
	SportsRegistrations *postgres.SportsRegistrationsModel
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

		var regs []*models.SportsRegistration
		var err error
		switch data.Type {
		case "currentAndUpcoming":
			regs, err = h.SportsRegistrations.SelectAllCurrentAndUpcoming(
				c.Request().Context())
		case "past":
			regs, err = h.SportsRegistrations.SelectAllPast(c.Request().Context())
		default:
			regs, err = h.SportsRegistrations.SelectAll(c.Request().Context())
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
