CREATE SCHEMA bank;

-- BALANCE TYPES
CREATE TABLE bank.balance_types (
  name VARCHAR(10) PRIMARY KEY);
INSERT INTO bank.balance_types (name)
VALUES ('Debit'), ('Credit');

-- ACCOUNT TYPES
CREATE TABLE bank.account_types (
  name VARCHAR(10) PRIMARY KEY,
  balance_type VARCHAR(10) NOT NULL
    REFERENCES bank.balance_types (name));
INSERT INTO bank.account_types
  (name, balance_type)
VALUES
  ('Chequing', 'Debit'),
  ('Savings', 'Debit'),
  ('Visa', 'Credit'),
  ('Mastercard', 'Credit');

-- CURRENCIES
CREATE TABLE bank.currencies (
  id VARCHAR(3) PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE);

-- BANKS
CREATE TABLE bank.banks (
  id VARCHAR(10) PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE);

-- ACCOUNTS
CREATE TABLE bank.accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL UNIQUE,
  bank_id VARCHAR(10) NOT NULL
    REFERENCES bank.banks (id),
  account_type VARCHAR(10) NOT NULL
    REFERENCES bank.account_types (name));

CREATE TABLE bank.account_payment_descriptions (
  account_id UUID NOT NULL
    REFERENCES bank.accounts (id),
  payment_type VARCHAR(10) NOT NULL,
  description VARCHAR(50) NOT NULL,
  PRIMARY KEY (account_id, payment_type, description));

-- TRANSACTIONS
CREATE TABLE bank.transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  date DATE NOT NULL,
  description VARCHAR(255) NOT NULL,
  description2 VARCHAR(255) NOT NULL DEFAULT '',
  debit DECIMAL(10, 2),
  credit DECIMAL(10, 2),
  currency_id VARCHAR(3) NOT NULL
    REFERENCES bank.currencies (id),
  account_number VARCHAR(255) NOT NULL DEFAULT '',
  card_number VARCHAR(255) NOT NULL DEFAULT '',
  cheque_number VARCHAR(255) NOT NULL DEFAULT '',
  account_id UUID NOT NULL
    REFERENCES bank.accounts (id));
