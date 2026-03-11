-- SPORTS
DROP TABLE IF EXISTS sports.registrations;
DROP SCHEMA IF EXISTS sports;

-- SCHOOL
DROP TABLE IF EXISTS school.reimbursement;
DROP TABLE IF EXISTS school.invoice;
DROP TABLE IF EXISTS school.school_year;
DROP TABLE IF EXISTS school.school;
DROP TABLE IF EXISTS school.grade;
DROP SCHEMA IF EXISTS school;

-- BOOK
DROP TABLE IF EXISTS book.invoice;

DROP TRIGGER IF EXISTS assign_unassigned
  ON book.journal_entry_account_entries;
DROP FUNCTION IF EXISTS book.assign_unassigned;
DROP TABLE IF EXISTS book.journal_entry_account_entries;
DROP TABLE IF EXISTS book.journal_entries;

DROP TRIGGER IF EXISTS scan_journal_entries_for_assigner_description_match
  ON book.assigner_bank_transaction_descriptions;
DROP FUNCTION IF EXISTS book.scan_journal_entries_for_assigner_description_match;
DROP TABLE IF EXISTS book.assigner_bank_transaction_descriptions;
DROP TABLE IF EXISTS book.assigners;

DROP TABLE IF EXISTS book.accounts;

DROP TABLE IF EXISTS book.currencies;

DROP TABLE IF EXISTS book.account_types;
DROP TABLE IF EXISTS book.balance_types;

DROP SCHEMA IF EXISTS book;

-- BANK
DROP TABLE IF EXISTS bank.transactions;
DROP TABLE IF EXISTS bank.account_payment_descriptions;
DROP TABLE IF EXISTS bank.accounts;
DROP TABLE IF EXISTS bank.banks;
DROP TABLE IF EXISTS bank.currencies;
DROP TABLE IF EXISTS bank.account_types;
DROP TABLE IF EXISTS bank.balance_types;
DROP SCHEMA IF EXISTS bank;
