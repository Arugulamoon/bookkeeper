package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/Arugulamoon/bookkeeper/pkg/config"
	"github.com/Arugulamoon/bookkeeper/pkg/models"
	"github.com/Arugulamoon/bookkeeper/pkg/models/postgres"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger

	DB *pgxpool.Pool

	bankTransactions *postgres.BankTransactionModel
	accounts         *postgres.AccountModel
	jEntries         *postgres.JournalEntryModel
	jAcctEntries     *postgres.JournalAccountEntryModel
}

func main() {
	var configFilename string
	flag.StringVar(&configFilename, "config", "", "path to config file")
	flag.Parse()

	if configFilename == "" {
		fmt.Println("missing config filename argument")
		os.Exit(1)
	}

	cfg, err := config.GetConfig(configFilename)
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

	ctx := context.Background()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,

		DB: db,

		bankTransactions: &postgres.BankTransactionModel{DB: db},
		accounts:         &postgres.AccountModel{DB: db},
		jEntries:         &postgres.JournalEntryModel{DB: db},
		jAcctEntries:     &postgres.JournalAccountEntryModel{DB: db},
	}

	app.loadCreditCardPaymentBankTransactions(ctx)
	app.loadOpaqueCreditCardPaymentBankTransactions(ctx)
	app.loadNonCreditCardPaymentBankTransactions(ctx)
}

// TODO: Rewrite queries with filtering using sql
func (app *application) loadCreditCardPaymentBankTransactions(
	ctx context.Context,
) {
	numInserted := 0

	alreadyImported, err := app.jAcctEntries.SelectAllAlreadyImportedBankTransactionIds(ctx)
	if err != nil {
		panic(err)
	}

	fetchedReceived, err := app.bankTransactions.SelectAllCreditCardPaymentsReceived(ctx)
	if err != nil {
		panic(err)
	}
	filteredReceived := filterAlreadyImportedBankTransactions(fetchedReceived, alreadyImported)

	fetchedPaid, err := app.bankTransactions.SelectAllPaymentsMadeToCreditCard(ctx)
	if err != nil {
		panic(err)
	}
	filteredPaid := filterAlreadyImportedBankTransactions(fetchedPaid, alreadyImported)

	for _, paid := range filteredPaid {
		for _, received := range filteredReceived {
			if paid.Credit == received.Debit {
				numInserted += app.loadCreditCardPaymentBankTransaction(
					ctx, received, paid)
			}
		}
	}

	fmt.Printf("Fetched %d payments received and %d payments made, filtered down to %d and %d and inserted %d bank transactions into journal\n",
		len(fetchedReceived), len(fetchedPaid), len(filteredReceived), len(filteredPaid), numInserted)
}

func (app *application) loadOpaqueCreditCardPaymentBankTransactions(
	ctx context.Context,
) {
	numInserted := 0

	alreadyImported, err := app.jAcctEntries.SelectAllAlreadyImportedBankTransactionIds(ctx)
	if err != nil {
		panic(err)
	}

	fetchedPaid, err := app.bankTransactions.SelectAllPaymentsMadeToOpaqueCreditCard(ctx)
	if err != nil {
		panic(err)
	}
	filteredPaid := filterAlreadyImportedBankTransactions(fetchedPaid, alreadyImported)

	for _, tx := range filteredPaid {
		numInserted += app.loadOpaqueCreditCardPaymentBankTransaction(ctx, tx)
	}

	fmt.Printf("Fetched %d Opaque CC payments made, filtered down to %d and inserted %d bank transactions into journal\n",
		len(fetchedPaid), len(filteredPaid), numInserted)
}

func (app *application) loadNonCreditCardPaymentBankTransactions(
	ctx context.Context,
) {
	numInserted := 0

	alreadyImported, err := app.jAcctEntries.SelectAllAlreadyImportedBankTransactionIds(ctx)
	if err != nil {
		panic(err)
	}

	fetched, err := app.bankTransactions.SelectAllNonCreditCardPayments(ctx)
	if err != nil {
		panic(err)
	}
	filtered := filterAlreadyImportedBankTransactions(fetched, alreadyImported)

	for _, tx := range filtered {
		numInserted += app.loadNonCreditCardPaymentBankTransaction(ctx, tx)
	}

	fmt.Printf("Fetched %d, filtered down to %d and inserted %d bank transactions into journal\n",
		len(fetched), len(filtered), numInserted)
}

func (app *application) loadCreditCardPaymentBankTransaction(
	ctx context.Context,
	receivedTx, paidTx *models.BankTransaction,
) int {
	debitAccount, err := app.accounts.SelectByBankAccountId(
		ctx, receivedTx.AccountId)
	if err != nil {
		panic(err)
	}
	creditAccount, err := app.accounts.SelectByBankAccountId(
		ctx, paidTx.AccountId)
	if err != nil {
		panic(err)
	}
	jEntry := JournalEntry{
		Date:        paidTx.Date, // use date from chequing acct
		Description: fmt.Sprintf("%s Payment from %s", debitAccount.Name, creditAccount.Name),
		Debit: &AccountEntry{
			AccountType:       debitAccount.AccountType,
			AccountName:       debitAccount.Name,
			Amount:            receivedTx.Debit,
			BankTransactionId: receivedTx.Id,
		},
		Credit: &AccountEntry{
			AccountType:       creditAccount.AccountType,
			AccountName:       creditAccount.Name,
			Amount:            paidTx.Credit,
			BankTransactionId: paidTx.Id,
		},
	}

	numInserted, err := app.insertJournalEntry(ctx, jEntry)
	if err != nil {
		panic(err)
	}
	return numInserted
}

func (app *application) loadOpaqueCreditCardPaymentBankTransaction(
	ctx context.Context,
	tx *models.BankTransaction,
) int {
	creditAccount, err := app.accounts.SelectByBankAccountId(ctx, tx.AccountId)
	if err != nil {
		panic(err)
	}
	jEntry := JournalEntry{
		Date:        tx.Date, // use date from chequing acct
		Description: fmt.Sprintf("Credit Card Payment from %s", creditAccount.Name),
		Debit: &AccountEntry{
			AccountType:       "Expense",
			AccountName:       "Unassigned",
			Amount:            tx.Credit,
			BankTransactionId: tx.Id,
		},
		Credit: &AccountEntry{
			AccountType:       creditAccount.AccountType,
			AccountName:       creditAccount.Name,
			Amount:            tx.Credit,
			BankTransactionId: tx.Id,
		},
	}

	numInserted, err := app.insertJournalEntry(ctx, jEntry)
	if err != nil {
		panic(err)
	}
	return numInserted
}

func (app *application) loadNonCreditCardPaymentBankTransaction(
	ctx context.Context,
	tx *models.BankTransaction,
) int {
	account, err := app.accounts.SelectByBankAccountId(ctx, tx.AccountId)
	if err != nil {
		panic(err)
	}

	jEntry := JournalEntry{
		Date:        tx.Date,
		Description: tx.Description,
	}

	if tx.Debit > 0.00 {
		amount := tx.Debit

		jEntry.Debit = &AccountEntry{
			AccountType:       account.AccountType,
			AccountName:       account.Name,
			Amount:            amount,
			BankTransactionId: tx.Id,
		}

		jEntry.Credit = &AccountEntry{
			AccountType:       "Revenue",
			AccountName:       "Unassigned",
			Amount:            amount,
			BankTransactionId: tx.Id,
		}

	} else if tx.Credit > 0.00 {
		amount := tx.Credit

		jEntry.Debit = &AccountEntry{
			AccountType:       "Expense",
			AccountName:       "Unassigned",
			Amount:            amount,
			BankTransactionId: tx.Id,
		}

		jEntry.Credit = &AccountEntry{
			AccountType:       account.AccountType,
			AccountName:       account.Name,
			Amount:            amount,
			BankTransactionId: tx.Id,
		}

	} else {
		log.Fatal("debit and credit cannot both be empty")
	}

	numInserted, err := app.insertJournalEntry(ctx, jEntry)
	if err != nil {
		panic(err)
	}
	return numInserted
}

func (app *application) insertJournalEntry(
	ctx context.Context,
	jEntry JournalEntry,
) (int, error) {
	jEntryId, err := app.jEntries.Insert(
		ctx,
		jEntry.Date, jEntry.Description)
	if err != nil {
		return 0, err
	}

	_, err = app.jAcctEntries.Insert(
		ctx,
		jEntryId, "Debit", nil,
		jEntry.Debit.AccountType, jEntry.Debit.AccountName, jEntry.Debit.Amount,
		jEntry.Debit.BankTransactionId)
	if err != nil {
		return 0, err
	}

	_, err = app.jAcctEntries.Insert(
		ctx,
		jEntryId, "Credit", nil,
		jEntry.Credit.AccountType, jEntry.Credit.AccountName, jEntry.Credit.Amount,
		jEntry.Credit.BankTransactionId)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

type JournalEntry struct {
	Date        time.Time
	Description string
	Debit       *AccountEntry
	Credit      *AccountEntry
}

type AccountEntry struct {
	AccountType       string
	AccountName       string
	Amount            float64
	BankTransactionId string
}

func filterAlreadyImportedBankTransactions(
	fetched []*models.BankTransaction,
	alreadyImported []string,
) []*models.BankTransaction {
	var filtered []*models.BankTransaction
	for _, bankTx := range fetched {
		if !slices.Contains(alreadyImported, bankTx.Id) {
			filtered = append(filtered, bankTx)
		}
	}
	return filtered
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Optional: Configure pool settings (e.g., max connections, lifetime)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse database config: %v", err)
	}
	config.MaxConns = 10
	config.MaxConnLifetime = 30 * time.Minute
	config.MinConns = 2

	// Establish the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}

	return pool, nil
}
