package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type AssignerHandler struct {
	DB *sql.DB
}

func (h *AssignerHandler) List() echo.HandlerFunc {
	return func(c *echo.Context) error {
		assignerModel := &postgres.AssignerModel{DB: h.DB}
		assgnrs, err := assignerModel.SelectAll()
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Render(http.StatusOK, "assigners.list.page.tmpl",
			map[string]any{
				"Assigners": assgnrs,
			})
	}
}

type CreateAssignerFormData struct {
	SearchedFor string `query:"search"`
}

func (h *AssignerHandler) CreateForm() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(CreateAssignerFormData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		accountModel := &postgres.AccountModel{DB: h.DB}
		accts, err := accountModel.SelectAll()
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
		entries, err := jAcctEntryModel.SelectAllByLikeDescription(
			data.SearchedFor)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Render(http.StatusOK, "assigners.create.page.tmpl",
			map[string]any{
				"SearchedFor":           data.SearchedFor,
				"Accounts":              accts,
				"JournalAccountEntries": entries,
			})
	}
}

type CreateAssignerData struct {
	Name       string `form:"name"`
	BankTxDesc string `form:"bankTxDesc"`
	Acct       string `form:"acct"`
}

func (h *AssignerHandler) Create() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(CreateAssignerData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		splittedAcct := strings.Split(data.Acct, "#") // AccountType#Name

		assignerModel := &postgres.AssignerModel{DB: h.DB}

		assignerId, err := assignerModel.Insert(
			data.Name, splittedAcct[0], splittedAcct[1])
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		_, err = assignerModel.InsertBankTransactionDescription(
			data.BankTxDesc, assignerId)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Redirect(http.StatusSeeOther, "/assigners")
	}
}

type ShowAssignerData struct {
	Id string `param:"id"`
}

func (h *AssignerHandler) Show() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(ShowAssignerData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		assignerModel := &postgres.AssignerModel{DB: h.DB}
		assgnr, err := assignerModel.Select(data.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
		entries, err := jAcctEntryModel.SelectAllByAssignerId(data.Id)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		url := c.Request().URL
		fromUrl := url.Path
		if url.RawQuery != "" {
			fromUrl += "?" + url.RawQuery
		}

		return c.Render(http.StatusOK, "assigners.show.page.tmpl",
			map[string]any{
				"Assigner":              assgnr,
				"JournalAccountEntries": entries,
				"FromUrl":               fromUrl,
			})
	}
}
