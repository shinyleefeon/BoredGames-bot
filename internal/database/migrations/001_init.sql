-- +goose Up
CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    guild_id TEXT NOT NULL,
    title TEXT NOT NULL,
    start_time DATETIMKE NOT NULL,
    reminder_sent BOOLEAN DEFAULT 0
);

CREATE TABLE USERS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    discord_id TEXT NOT NULL,
    username TEXT NOT NULL,
    streak INTEGER DEFAULT 0,
    last_victory_time TEXT
);

CREATE TABLE participants (
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    time_since_last_victory TEXT,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE participants;
DROP TABLE users;
DROP TABLE events;

