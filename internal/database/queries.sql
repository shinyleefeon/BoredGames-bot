-- name: CreateEvent :one
INSERT INTO events (discord_event_id, guild_id, title, start_time, reminder_sent) VALUES (?, ?, ?, ?, 0)
RETURNING id, discord_event_id, guild_id, title, start_time, reminder_sent;

-- name: AddParticipant :exec
INSERT OR IGNORE INTO participants (event_id, user_id) VALUES (?, ?);

-- name: GetEventsForReminder :many
-- Finds events starting within the next 30 minutes (approx) that haven't been sent
SELECT id, discord_event_id, guild_id, title, start_time
FROM events 
WHERE start_time <= datetime('now', '+30 minutes') 
AND reminder_sent = 0;

-- name: MarkReminderSent :exec
UPDATE events 
SET reminder_sent = 1 
WHERE id = ?;

-- name: GetParticipants :many
SELECT users.discord_id
FROM participants
JOIN users ON users.id = participants.user_id
WHERE participants.event_id = ?;

-- name: GetUserByDiscordID :one
SELECT id, discord_id, username, streak, last_victory_time 
FROM users 
WHERE discord_id = ?;

-- name: CreateUser :one
INSERT INTO users (discord_id, username) VALUES (?, ?)
RETURNING id, discord_id, username, streak, last_victory_time;

-- name: UpdateUserVictory :exec
UPDATE users 
SET streak = streak + 1, last_victory_time = datetime('now') 
WHERE id = ?;

-- name: GetUserStreak :one
SELECT streak 
FROM users 
WHERE id = ?;

-- name: GetUserLastVictory :one
SELECT last_victory_time 
FROM users 
WHERE id = ?;

-- name: IncrementStreak :exec
UPDATE users
SET streak = streak + 1, last_victory_time = datetime('now')
WHERE id = ?;

-- name: resetStreak :exec
UPDATE users
SET streak = 0
WHERE id = ?;

-- name: GetEventByEventID :one
SELECT id, discord_event_id, guild_id, title, start_time, reminder_sent
FROM events
WHERE discord_event_id = ?;