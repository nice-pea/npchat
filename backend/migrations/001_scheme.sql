CREATE TABLE roles
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    permissions TEXT
);

CREATE TABLE localization
(
    category TEXT NOT NULL,
    code     TEXT NOT NULL,
    locale   TEXT NOT NULL,
    text     TEXT NOT NULL
);
