package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type JournalAccountEntryHandler struct {
	DB *sql.DB
}

type EditJournalAccountEntryFormData struct {
	Id      string `param:"id"`
	FromUrl string `query:"fromUrl"`
}

func (h *JournalAccountEntryHandler) EditForm() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(EditJournalAccountEntryFormData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
		entry, err := jAcctEntryModel.Select(data.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		accountModel := &postgres.AccountModel{DB: h.DB}
		accts, err := accountModel.SelectAll()
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Render(http.StatusOK, "jacctentries.edit.page.tmpl",
			map[string]any{
				"FromUrl":             data.FromUrl,
				"JournalAccountEntry": entry,
				"Accounts":            accts,
			})
	}
}

type EditJournalAccountEntryData struct {
	Id      string `param:"id"`
	FromUrl string `query:"fromUrl"`
	Acct    string `form:"acct"`
}

func (h *JournalAccountEntryHandler) Edit() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(EditJournalAccountEntryData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		acctParts := strings.Split(data.Acct, "#") // AccountType#Name
		acctType := acctParts[0]
		acctName := acctParts[1]

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
		_, err := jAcctEntryModel.UpdateAccount(data.Id, acctType, acctName)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Redirect(http.StatusSeeOther, data.FromUrl)
	}
}

type SplitJournalAccountEntryFormData struct {
	Id      string `param:"id"`
	FromUrl string `query:"fromUrl"`
}

func (h *JournalAccountEntryHandler) SplitForm() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(SplitJournalAccountEntryFormData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
		entry, err := jAcctEntryModel.Select(data.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		accountModel := &postgres.AccountModel{DB: h.DB}
		accts, err := accountModel.SelectAll()
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Render(http.StatusOK, "jacctentries.split.page.tmpl",
			map[string]any{
				"FromUrl":             data.FromUrl,
				"JournalAccountEntry": entry,
				"Accounts":            accts,
			})
	}
}

type SplitJournalAccountEntryData struct {
	Id      string  `param:"id"`
	FromUrl string  `query:"fromUrl"`
	Amount1 float64 `form:"amount1"`
	Amount2 float64 `form:"amount2"`
	Acct2   string  `form:"acct2"`
}

func (h *JournalAccountEntryHandler) Split() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(SplitJournalAccountEntryData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}

		entry, err := jAcctEntryModel.Select(data.Id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return c.String(http.StatusNotFound, "not found")
			} else {
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}
		}

		// TODO: Replace money in entire app as cents integer
		if data.Amount1*100+data.Amount2*100 != entry.Amount*100 {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("amounts must equal $%.2f", entry.Amount))
		}

		acctParts := strings.Split(data.Acct2, "#") // AccountType#Name
		acct2AcctType := acctParts[0]
		acct2Name := acctParts[1]

		_, err = jAcctEntryModel.UpdateAmount(data.Id, data.Amount1)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		_, err = jAcctEntryModel.Insert(
			entry.JournalEntryId,
			entry.BalanceType,
			entry.AssignerId,
			acct2AcctType, acct2Name,
			data.Amount2, entry.BankTransactionId,
		)
		if err != nil {
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		return c.Redirect(http.StatusSeeOther, data.FromUrl)
	}
}
