package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type SportsMembershipHandler struct {
	SportsMemberships *postgres.SportsMembershipModel
}

func (h *SportsMembershipHandler) List() echo.HandlerFunc {
	return func(c *echo.Context) error {
		memberships, err := h.SportsMemberships.SelectAll(
			c.Request().Context())
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

		membership, err := h.SportsMemberships.Select(
			c.Request().Context(),
			data.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		games, err := h.SportsMemberships.SelectAllGames(
			c.Request().Context(),
			data.Id)
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
