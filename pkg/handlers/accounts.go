package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type AccountHandler struct {
	DB *sql.DB
}

func (h *AccountHandler) List() echo.HandlerFunc {
	return func(c *echo.Context) error {
		accountModel := &postgres.AccountModel{DB: h.DB}
		accts, err := accountModel.SelectAll()
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Render(http.StatusOK, "accounts.list.page.tmpl",
			map[string]any{
				"Accounts": accts,
			})
	}
}

func (h *AccountHandler) CreateForm() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.Render(http.StatusOK, "accounts.create.page.tmpl",
			map[string]any{})
	}
}

type CreateAccountFormData struct {
	AcctType  string `form:"acctType"`
	Name      string `form:"name"`
	SortOrder int    `form:"sortOrder"`
}

func (h *AccountHandler) Create() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(CreateAccountFormData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if data.SortOrder <= 0 {
			data.SortOrder = 1000
		}

		accountModel := &postgres.AccountModel{DB: h.DB}
		_, err := accountModel.Insert(
			data.AcctType, data.Name, nil, data.SortOrder)
		if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther,
			fmt.Sprintf("/account/%s/%s", data.AcctType, data.Name))
	}
}

type ShowAccountData struct {
	AcctType string `param:"acctType"`
	AcctName string `param:"acctName"`
}

func (h *AccountHandler) Show() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(ShowAccountData)
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
		entries, err := jAcctEntryModel.SelectAllByAccountId(
			data.AcctType, data.AcctName)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		url := c.Request().URL
		fromUrl := url.Path
		if url.RawQuery != "" {
			fromUrl += "?" + url.RawQuery
		}

		return c.Render(http.StatusOK, "accounts.show.page.tmpl",
			map[string]any{
				"Account":               acct,
				"JournalAccountEntries": entries,
				"FromUrl":               fromUrl,
			})
	}
}
