CREATE TABLE roles
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    permissions TEXT
);

CREATE TABLE localizations
(
    category TEXT NOT NULL,
    code     TEXT NOT NULL,
    locale   TEXT NOT NULL,
    text     TEXT NOT NULL
);

CREATE TABLE users
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL
);

CREATE TABLE users
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL
);

CREATE TABLE credentials
(
    user_id INTEGER,
    key     TEXT NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE
);

CREATE TABLE sessions
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER REFERENCES users ON DELETE CASCADE,
    token      TEXT     NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
);