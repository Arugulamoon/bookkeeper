package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type AnnualReportHandler struct {
	DB *sql.DB
}

type ShowAnnualReportData struct {
	Year int `query:"year"`
}

func (h *AnnualReportHandler) Show() echo.HandlerFunc {
	return func(c *echo.Context) error {
		data := new(ShowAnnualReportData)
		if err := c.Bind(data); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		var year int
		if data.Year > 0 {
			year = data.Year
		} else {
			year = time.Now().Year()
		}

		accountModel := &postgres.AccountModel{DB: h.DB}
		accts, err := accountModel.SelectAll()
		if err != nil {
			c.Logger().Error(err.Error())
			return c.String(http.StatusInternalServerError,
				"internal server error")
		}

		acctTotals := map[string][]*models.AccountAnnualTotal{
			"Asset":     make([]*models.AccountAnnualTotal, 0),
			"Liability": make([]*models.AccountAnnualTotal, 0),
			"Revenue":   make([]*models.AccountAnnualTotal, 0),
			"Expense":   make([]*models.AccountAnnualTotal, 0),
		}
		acctTypeTotals := map[string]*models.AccountTypeAnnualTotal{
			"Asset": {
				MonthlyTotals: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				AnnualTotal:   0,
			},
			"Liability": {
				MonthlyTotals: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				AnnualTotal:   0,
			},
			"Revenue": {
				MonthlyTotals: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				AnnualTotal:   0,
			},
			"Expense": {
				MonthlyTotals: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				AnnualTotal:   0,
			},
		}
		for _, acct := range accts {
			jAcctEntryModel := &postgres.JournalAccountEntryModel{DB: h.DB}
			entries, err := jAcctEntryModel.SelectAllByAccountForYear(
				acct.AccountType, acct.Name, year)
			if err != nil {
				c.Logger().Error(err.Error())
				return c.String(http.StatusInternalServerError,
					"internal server error")
			}

			if len(entries) > 0 {
				acctTotal := &models.AccountAnnualTotal{
					Account:       acct,
					MonthlyTotals: []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // jan, feb etc
					AnnualTotal:   0,
				}
				for _, entry := range entries {
					monthIdx := int(entry.Date.Month()) - 1
					if entry.AccountBalanceType == entry.BalanceType {
						acctTotal.MonthlyTotals[monthIdx] += entry.Amount
						acctTotal.AnnualTotal += entry.Amount
						acctTypeTotals[acct.AccountType].MonthlyTotals[monthIdx] += entry.Amount
						acctTypeTotals[acct.AccountType].AnnualTotal += entry.Amount
					} else {
						acctTotal.MonthlyTotals[monthIdx] -= entry.Amount
						acctTotal.AnnualTotal -= entry.Amount
						acctTypeTotals[acct.AccountType].MonthlyTotals[monthIdx] -= entry.Amount
						acctTypeTotals[acct.AccountType].AnnualTotal -= entry.Amount
					}
				}
				acctTotals[acct.AccountType] = append(acctTotals[acct.AccountType], acctTotal)
			}
		}

		return c.Render(http.StatusOK, "annuals.page.tmpl",
			map[string]any{
				"AnnualReportYear":         year,
				"AccountMonthlyTotals":     acctTotals,
				"AccountTypeMonthlyTotals": acctTypeTotals,
			})
	}
}
