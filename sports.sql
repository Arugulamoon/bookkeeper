CREATE SCHEMA sports;

CREATE TABLE sports.registrations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  price INTEGER NOT NULL DEFAULT 0,
  regular_price INTEGER,
  discount INTEGER,
  tax INTEGER,
  location VARCHAR(255),
  day VARCHAR(50),
  start_time VARCHAR(5),
  start_time_range VARCHAR(5),
  end_time_range VARCHAR(5),
  duration INTEGER,
  start_date DATE,
  end_date DATE,
  sessions INTEGER
);

CREATE TABLE sports.memberships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  season_year VARCHAR(7) NOT NULL,
  season_type VARCHAR(50) NOT NULL,
  location VARCHAR(255) NOT NULL
);

CREATE TABLE sports.membership_games (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  membership_id UUID NOT NULL
    REFERENCES sports.memberships (id),
  date DATE NOT NULL,
  start_time VARCHAR(5) NOT NULL,
  opponent VARCHAR(50) NOT NULL,
  notes VARCHAR(255),
  location VARCHAR(255),
  event_id CHAR(26)
);
