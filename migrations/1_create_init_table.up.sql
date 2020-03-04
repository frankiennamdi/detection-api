CREATE TABLE events (
    uuid TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    timestamp NUMERIC NOT NULL,
    ip  TEXT NOT NULL
);

CREATE UNIQUE INDEX events_username_timestamp_unq ON events(username, timestamp);