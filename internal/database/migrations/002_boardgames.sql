-- +goose Up
CREATE TABLE boardgames (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    category TEXT,
    min_players INTEGER,
    max_players INTEGER,
    play_time INTEGER,
    description TEXT,
    previous_winner TEXT,
    played_yet BOOLEAN,
    liked_it BOOLEAN,
    rules_link TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_player_count ON boardgames (min_players, max_players);


-- +goose Down
DROP TABLE boardgames;