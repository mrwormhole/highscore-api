CREATE TABLE highscores (
  id BIGSERIAL PRIMARY KEY,
  username text NOT NULL,
  score int NOT NULL
);