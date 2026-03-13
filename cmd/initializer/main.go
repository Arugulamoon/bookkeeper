package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"

	"github.com/Arugulamoon/bookkeeper/pkg/config"
	"github.com/Arugulamoon/bookkeeper/pkg/google"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
	"github.com/Arugulamoon/bookkeeper/pkg/yamlmodels"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	db *sql.DB

	bankCurrencies   *postgres.BankCurrencyModel
	banks            *postgres.BankModel
	bankAccounts     *postgres.BankAccountModel
	bankTransactions *postgres.BankTransactionModel

	sportsMemberships *postgres.SportsMembershipModel

	calendar *google.EventingGCalendar
}

func main() {
	var configFilename, dataFilename string
	flag.StringVar(&configFilename, "config", "", "path to config file")
	flag.StringVar(&dataFilename, "data", "./testdata/data.yaml", "path to data file")
	flag.Parse()

	if configFilename == "" {
		fmt.Println("missing config filename argument")
		os.Exit(1)
	}

	cfg, err := config.GetConfig(configFilename)
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(dataFilename)
	if err != nil {
		panic(err)
	}

	var data yamlmodels.ConfigData
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		panic(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.DSN())
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	gcalsvc, err := calendar.NewService(context.Background(),
		option.WithHTTPClient(google.GetClient(cfg.Google.Auth.Dir)))
	if err != nil {
		errorLog.Fatalf("Unable to create Calendar service: %v", err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,

		db: db,

		bankCurrencies:   &postgres.BankCurrencyModel{DB: db},
		banks:            &postgres.BankModel{DB: db},
		bankAccounts:     &postgres.BankAccountModel{DB: db},
		bankTransactions: &postgres.BankTransactionModel{DB: db},

		sportsMemberships: &postgres.SportsMembershipModel{DB: db},

		calendar: &google.EventingGCalendar{
			Service:   gcalsvc.Events,
			Calendars: cfg.Google.GetCalendars(),
		},
	}

	app.initBank(data.Bank)
	app.initBook(data.Book)
	app.initSports(data.Sports)
	app.initSchool(data.School)
}

func (app *application) initBank(data yamlmodels.BankData) {
	for _, currency := range data.Currencies {
		_, err := app.bankCurrencies.Insert(currency.Id, currency.Name)
		if err != nil {
			panic(err)
		}
	}

	for _, bank := range data.Banks {
		err := app.banks.Insert(bank.Id, bank.Name)
		if err != nil {
			panic(err)
		}
		for _, acct := range bank.Accounts {
			acctId, err := app.bankAccounts.Insert(acct.Name, bank.Id, acct.Type)
			if err != nil {
				panic(err)
			}

			for _, filename := range acct.Files {
				file, err := os.Open(filename)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				// Read the file content into a byte slice
				data, err := io.ReadAll(file)
				if err != nil {
					panic(err)
				}

				// Trim UTF-8 BOM prefix
				bom := []byte{0xEF, 0xBB, 0xBF} // UTF-8 BOM
				sanitizedData := bytes.TrimPrefix(data, bom)

				reader := csv.NewReader(bytes.NewReader(sanitizedData))

				rows, err := reader.ReadAll()
				if err != nil {
					panic(err)
				}

				numInsertedTxs := 0
				switch bank.Id {
				case "RBC":
					numInsertedTxs = app.processRBCTransactions(acctId, rows)
				case "CIBC":
					numInsertedTxs = app.processCIBCTransactions(acctId, rows)
				default:
					panic("unknown bank")
				}

				fmt.Printf("%s: Inserted %d transactions from %s\n",
					acct.Name, numInsertedTxs, filename)
			}

			for _, desc := range acct.PaymentRecdDescs {
				err := app.bankAccounts.InsertPaymentDescription(
					acctId, "Received", desc)
				if err != nil {
					panic(err)
				}
			}

			for _, desc := range acct.PaymentMadeDescs {
				err := app.bankAccounts.InsertPaymentDescription(
					acctId, "Made", desc)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (app *application) processRBCTransactions(
	acctId string, rows [][]string,
) int {
	rawTxs := rows[1:] // headers in first row; skip them
	numToInsert := len(rawTxs)
	numInserted := 0

	for _, rawTx := range rawTxs {
		tx := NewRBCTransaction(rawTx)

		debit := 0.00
		credit := 0.00
		var currency string
		if tx.CAD != 0.00 {
			amount := math.Abs(tx.CAD)
			if tx.CAD > 0.00 {
				debit = amount
			} else {
				credit = amount
			}
			currency = "CAD"
		} else if tx.USD != 0.00 {
			amount := math.Abs(tx.USD)
			if tx.USD > 0.00 {
				debit = amount
			} else {
				credit = amount
			}
			currency = "USD"
		} else {
			panic("CAD and USD cannot both be zero")
		}

		id, err := app.bankTransactions.InsertRBC(
			tx.TransactionDate,
			tx.Description, tx.Description2,
			debit, credit, currency,
			tx.AccountNumber, tx.ChequeNumber,
			acctId)
		if err != nil {
			panic(err)
		}
		if id != "" {
			numInserted++
		}
	}

	if numInserted != numToInsert {
		panic(fmt.Sprintf("inserted (%d) did not match num to insert (%d)\n",
			numInserted, numToInsert))
	}

	return numInserted
}

type RBCTransaction struct {
	AccountType, AccountNumber string
	TransactionDate            time.Time
	ChequeNumber               string
	Description, Description2  string
	CAD, USD                   float64
}

func NewRBCTransaction(row []string) *RBCTransaction {
	txDate, err := time.Parse("1/2/2006", row[2])
	if err != nil {
		panic(err)
	}

	return &RBCTransaction{
		AccountType:     row[0],
		AccountNumber:   row[1],
		TransactionDate: txDate,
		ChequeNumber:    row[3],
		Description:     row[4],
		Description2:    row[5],
		CAD:             castToFloat(row[6]),
		USD:             castToFloat(row[7]),
	}
}

func (app *application) processCIBCTransactions(
	acctId string, rows [][]string,
) int {
	rawTxs := rows // headers NOT in first row
	numToInsert := len(rawTxs)
	numInserted := 0

	for _, rawTx := range rawTxs {
		tx := NewCIBCTransaction(rawTx)
		id, err := app.bankTransactions.InsertCIBC(
			tx.Date,
			tx.Description,
			tx.Debit, tx.Credit,
			tx.CardNumber,
			acctId)
		if err != nil {
			panic(err)
		}
		if id != "" {
			numInserted++
		}
	}

	if numInserted != numToInsert {
		panic(fmt.Sprintf("inserted (%d) did not match num to insert (%d)\n",
			numInserted, numToInsert))
	}

	return numInserted
}

type CIBCTransaction struct {
	Date          time.Time
	Description   string
	Credit, Debit float64
	CardNumber    string
}

func NewCIBCTransaction(row []string) *CIBCTransaction {
	date, err := time.Parse("2006-01-02", row[0])
	if err != nil {
		panic(err)
	}

	return &CIBCTransaction{
		Date:        date,
		Description: row[1],
		Credit:      castToFloat(row[2]),
		Debit:       castToFloat(row[3]),
		CardNumber:  row[4],
	}
}

func (app *application) initBook(data yamlmodels.BookData) {
	for _, currency := range data.Currencies {
		currencyModel := &postgres.BookCurrencyModel{DB: app.db}
		_, err := currencyModel.Insert(currency.Id, currency.Name)
		if err != nil {
			panic(err)
		}
	}

	for _, acct := range data.Accounts {
		var bankAccountId *string
		if acct.BankAccount != nil {
			id, err := app.bankAccounts.GetId(*acct.BankAccount)
			if err != nil {
				panic(err)
			}
			if id != nil {
				bankAccountId = id
			}
		}

		sortOrder := 1000
		if acct.SortOrder != nil {
			sortOrder = *acct.SortOrder
		}

		accountModel := &postgres.AccountModel{DB: app.db}
		_, err := accountModel.Insert(
			acct.AccountType, acct.Name, bankAccountId, sortOrder,
		)
		if err != nil {
			panic(err)
		}
	}

	for _, assigner := range data.Assigners {
		assignerModel := &postgres.AssignerModel{DB: app.db}

		id, err := assignerModel.Insert(
			assigner.Name, assigner.AccountType, assigner.Account)
		if err != nil {
			panic(err)
		}
		for _, desc := range assigner.Descriptions {
			_, err := assignerModel.InsertBankTransactionDescription(
				desc, id)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (app *application) initSports(data yamlmodels.SportsData) {
	for _, reg := range data.Registrations {
		sportsRegModel := &postgres.SportsRegistrationsModel{DB: app.db}
		_, err := sportsRegModel.Insert(
			reg.Name,
			reg.Price.Total,
			reg.Price.Regular,
			reg.Price.Discount,
			reg.Price.Tax,
			reg.Location,
			reg.Time.Day,
			reg.Time.Start,
			reg.Time.Range.Start,
			reg.Time.Range.End,
			reg.Time.Duration,
			reg.Date.Start,
			reg.Date.End,
			reg.Sessions,
		)
		if err != nil {
			panic(err)
		}
	}

	app.initSportsMemberships(data.Memberships)
}

func (app *application) initSportsMemberships(
	memberships []yamlmodels.Membership,
) {
	for _, membership := range memberships {
		membershipId, err := app.sportsMemberships.Insert(membership.Name,
			membership.Season.Year, membership.Season.Type, membership.Location)
		if err != nil {
			panic(err)
		}
		for _, game := range membership.Games {
			gameId, err := app.sportsMemberships.InsertGame(membershipId,
				game.Date, game.Time.Start, game.Opponent, game.Notes, game.Location,
				game.Event.Id)
			if err != nil {
				panic(err)
			}
			if game.Event.Id == "" {
				// TODO: Move into function/method
				summary := fmt.Sprintf("🚨🏒🥅 %s at Ottawa Charge",
					game.Opponent)
				if game.Notes != "" {
					summary += fmt.Sprintf(" (%s)", game.Notes)
				}

				var location string
				if game.Location != "" {
					location = game.Location
				} else {
					location = membership.Location
				}

				var eventId string
				if game.Time.Start == "" || game.Time.Start == "TBD" {
					eventId = app.calendar.CreateAllDayEvent(membership.Calendar,
						game.Date.Format(time.DateOnly), summary, location)
				} else {
					eventId = app.calendar.CreateEvent(membership.Calendar,
						game.Date.Format(time.DateOnly), game.Time.Start, summary, location)
				}

				_, err = app.sportsMemberships.UpdateGameEventId(gameId, eventId)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (app *application) initSchool(data yamlmodels.SchoolData) {
	schoolsModel := &postgres.SchoolModel{DB: app.db}
	for _, grade := range data.Grades {
		err := schoolsModel.InsertGrade(grade.Id, grade.Name)
		if err != nil {
			panic(err)
		}
	}

	for _, school := range data.Schools {
		err := schoolsModel.InsertSchool(school.Id, school.Name,
			school.Address, school.Phone, school.Principal)
		if err != nil {
			panic(err)
		}
	}

	for _, schoolYear := range data.SchoolYears {
		err := schoolsModel.InsertSchoolYear(schoolYear.Year,
			schoolYear.SchoolId, schoolYear.GradeId,
			schoolYear.Teacher, schoolYear.Education)
		if err != nil {
			panic(err)
		}
	}

	for _, inv := range data.Invoices {
		invoicesModel := &postgres.InvoicesModel{DB: app.db}
		invId, err := invoicesModel.Insert(inv.DueDate, inv.Description, inv.Amount)
		if err != nil {
			panic(err)
		}
		schoolExpensesModel := &postgres.SchoolExpensesModel{DB: app.db}
		schoolId, err := schoolExpensesModel.InsertInvoice(
			invId, inv.SchoolYear, inv.School, inv.Grade,
			inv.Event.Id, inv.DatePaid, inv.EventMarkedPaid)
		if err != nil {
			panic(err)
		}
		if inv.Reimbursement != nil {
			_, err = schoolExpensesModel.InsertReimbursement(
				schoolId,
				inv.Reimbursement.Split,
				inv.Reimbursement.Amount,
				inv.Reimbursement.Date)
			if err != nil {
				panic(err)
			}
		}
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func castToFloat(s string) float64 {
	if s == "" {
		return 0.00
	}
	float, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return float
}
