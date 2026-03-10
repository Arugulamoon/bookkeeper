CREATE SCHEMA book;

-- BALANCE TYPES
CREATE TABLE book.balance_types (
  name VARCHAR(10) PRIMARY KEY);
INSERT INTO book.balance_types (name)
VALUES ('Debit'), ('Credit');

-- ACCOUNT TYPES
CREATE TABLE book.account_types (
  name VARCHAR(10) PRIMARY KEY,
  balance_type VARCHAR(10) NOT NULL
    REFERENCES book.balance_types (name),
  sort_order INTEGER NOT NULL DEFAULT 100);
INSERT INTO book.account_types
  (name, balance_type, sort_order)
VALUES
  ('Asset', 'Debit', 1),
  ('Liability', 'Credit', 2),
  ('Revenue', 'Credit', 3),
  ('Expense', 'Debit', 4);

-- CURRENCIES
CREATE TABLE book.currencies (
  id VARCHAR(3) PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE
);

-- ACCOUNTS
CREATE TABLE book.accounts (
  account_type VARCHAR(10) NOT NULL
    REFERENCES book.account_types (name),
  name VARCHAR(255) NOT NULL,
  bank_account_id UUID,
  sort_order INTEGER NOT NULL DEFAULT 1000,
  PRIMARY KEY (account_type, name)
);

-- ASSIGNERS
CREATE TABLE book.assigners (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) UNIQUE NOT NULL,
  account_type VARCHAR(10) NOT NULL,
  account_name VARCHAR(255) NOT NULL,
  FOREIGN KEY (account_type, account_name)
    REFERENCES book.accounts (account_type, name) ON UPDATE CASCADE
);

-- matching criteria
-- TODO (10): add other matching tables for diff desc matches, amount ranges etc
CREATE TABLE book.assigner_bank_transaction_descriptions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  bank_transaction_description VARCHAR(255) UNIQUE NOT NULL,
  assigner_id UUID NOT NULL
    REFERENCES book.assigners (id)
);

CREATE OR REPLACE FUNCTION book.scan_journal_entries_for_assigner_description_match()
RETURNS TRIGGER
AS $$
DECLARE
  accttype VARCHAR(10);
  acctname VARCHAR(255);
BEGIN
  SELECT account_type, account_name
  INTO accttype, acctname
  FROM book.assigners
  WHERE id = NEW.assigner_id;

  UPDATE book.journal_entry_account_entries AS acctentry
  SET
    assigner_id = NEW.assigner_id,
    account_type = accttype,
    account_name = acctname
  FROM
    book.journal_entries AS jentry
  WHERE
    acctentry.journal_entry_id = jentry.id AND
    acctentry.account_type IN ('Revenue', 'Expense') AND
    jentry.description ILIKE CONCAT('%', NEW.bank_transaction_description, '%');

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER scan_journal_entries_for_assigner_description_match
AFTER INSERT ON book.assigner_bank_transaction_descriptions
FOR EACH ROW
EXECUTE FUNCTION book.scan_journal_entries_for_assigner_description_match();

-- JOURNAL
CREATE TABLE book.journal_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  date DATE NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE book.journal_entry_account_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  journal_entry_id UUID
    REFERENCES book.journal_entries (id),
  balance_type VARCHAR(10) NOT NULL
    REFERENCES book.balance_types (name),
  assigner_id UUID
    REFERENCES book.assigners (id),
  account_type VARCHAR(10) NOT NULL,
  account_name VARCHAR(255) NOT NULL,
  amount DECIMAL(10, 2) NOT NULL,
  bank_transaction_id UUID,
  FOREIGN KEY (account_type, account_name)
    REFERENCES book.accounts (account_type, name) ON UPDATE CASCADE
);

CREATE OR REPLACE FUNCTION book.assign_unassigned()
RETURNS TRIGGER
AS $$
DECLARE
  jentrydesc VARCHAR(255);
  assignerid UUID;
  accttype VARCHAR(10);
  acctname VARCHAR(255);
BEGIN
  IF NEW.account_type IN ('Revenue', 'Expense') AND NEW.account_name = 'Unassigned' THEN
    SELECT description
    INTO jentrydesc
    FROM book.journal_entries
    WHERE id = NEW.journal_entry_id;

    SELECT assigners.id, assigners.account_type, assigners.account_name
    INTO assignerid, accttype, acctname
  	FROM book.assigner_bank_transaction_descriptions AS expdescs
    INNER JOIN book.assigners AS assigners
      ON expdescs.assigner_id = assigners.id
  	WHERE jentrydesc ILIKE CONCAT('%', expdescs.bank_transaction_description, '%');

    IF assignerid IS NOT NULL AND accttype IS NOT NULL AND acctname IS NOT NULL THEN
  		UPDATE book.journal_entry_account_entries
  		SET
        assigner_id = assignerid,
        account_type = accttype,
        account_name = acctname
  		WHERE id = NEW.id;
    END IF;
  END IF;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER assign_unassigned
AFTER INSERT ON book.journal_entry_account_entries
FOR EACH ROW
EXECUTE FUNCTION book.assign_unassigned();

-- TODO: Wire into school.invoice
CREATE TABLE book.invoice (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  due_date DATE NOT NULL,
  description VARCHAR(255) NOT NULL,
  amount INTEGER NOT NULL);
-- TODO: Add event id etc to schema notifications
