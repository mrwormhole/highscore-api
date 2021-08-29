CREATE TABLE highscores (
  id BIGSERIAL PRIMARY KEY,
  username text NOT NULL UNIQUE,
  score int NOT NULL
);