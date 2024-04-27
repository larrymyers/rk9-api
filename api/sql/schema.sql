CREATE TABLE
  events (
    id text PRIMARY key,
    name text,
    location text,
    start_date DATE,
    end_date DATE,
    url text
  );

CREATE TABLE
  players (
    id text PRIMARY key,
    name text,
    country text,
    wins INTEGER,
    losses INTEGER,
    ties INTEGER,
    points INTEGER,
    decklist_url text,
    standing INTEGER
  );

CREATE TABLE
  matches (
    id text PRIMARY key,
    pod INTEGER NOT NULL,
    round_number INTEGER NOT NULL,
    table_number text,
    player1_id text,
    player2_id text,
    winner_id text,
    is_tie BOOLEAN DEFAULT FALSE
  );