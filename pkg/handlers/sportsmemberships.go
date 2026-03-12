package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
	"github.com/labstack/echo/v5"
)

type SportsMembershipHandler struct {
	DB *sql.DB
}

func (h *SportsMembershipHandler) List() echo.HandlerFunc {
	return func(c *echo.Context) error {
		m := &postgres.SportsMembershipModel{DB: h.DB}
		memberships, err := m.SelectAll()
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "sportsmemberships.list.page.tmpl",
			map[string]any{
				"SportsMemberships": memberships,
			})
	}
}

type ShowSportsMembershipData struct {
	Id string `param:"id"`
}

func (h *SportsMembershipHandler) Show() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(ShowSportsMembershipData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		m := &postgres.SportsMembershipModel{DB: h.DB}

		membership, err := m.Select(data.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		games, err := m.SelectAllGames(data.Id)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "sportsmemberships.show.page.tmpl",
			map[string]any{
				"Membership": membership,
				"Games":      games,
			})
	}
}
