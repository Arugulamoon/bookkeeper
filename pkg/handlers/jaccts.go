package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type JournalAccountHandler struct {
	DB *sql.DB
}

type ShowJournalAccountData struct {
	AcctType string `param:"acctType"`
	AcctName string `param:"acctName"`
	Year     int    `query:"year"`
	Month    int    `query:"month"`
}

func (h *JournalAccountHandler) Show() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(ShowJournalAccountData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		accountModel := &postgres.AccountModel{DB: h.DB}
		acct, err := accountModel.Select(data.AcctType, data.AcctName)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
		entries, err := jAcctEntryModel.SelectAllByAccountForMonth(
			data.AcctType, data.AcctName, data.Year, data.Month)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		url := c.Request().URL
		fromUrl := url.Path
		if url.RawQuery != "" {
			fromUrl += "?" + url.RawQuery
		}

		return c.Render(http.StatusOK, "jaccts.show.page.tmpl",
			map[string]any{
				"Account":               acct,
				"JournalAccountEntries": entries,
				"FromUrl":               fromUrl,
			})
	}
}
