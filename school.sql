CREATE SCHEMA school;

CREATE TABLE school.grade (
  id VARCHAR(2) PRIMARY KEY,
  name VARCHAR(20) NOT NULL);

CREATE TABLE school.school (
  id VARCHAR(10) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  address VARCHAR(255),
  phone CHAR(10),
  principal VARCHAR(50));

CREATE TABLE school.school_year (
  school_year VARCHAR(7) NOT NULL,
  school_id VARCHAR(10) NOT NULL
    REFERENCES school.school (id),
  grade_id VARCHAR(2) NOT NULL
    REFERENCES school.grade (id),
  teacher VARCHAR(50),
  education VARCHAR(255),
  PRIMARY KEY (school_year, school_id, grade_id));

CREATE TABLE school.invoice (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  invoice_id UUID NOT NULL
    REFERENCES book.invoice (id),
  school_year VARCHAR(7) NOT NULL,
  school_id VARCHAR(10) NOT NULL,
  grade_id VARCHAR(2) NOT NULL,
  event_id CHAR(26),
  date_paid DATE,
  event_marked_paid BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (school_year, school_id, grade_id)
    REFERENCES school.school_year (school_year, school_id, grade_id)
    ON UPDATE CASCADE);

CREATE TABLE school.reimbursement (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  invoice_id UUID NOT NULL
    REFERENCES school.invoice (id),
  split CHAR(5),
  amount INTEGER,
  date DATE);
