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
