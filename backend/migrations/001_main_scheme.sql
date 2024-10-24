CREATE TABLE roles
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    chat_id     INTEGER REFERENCES chats ON DELETE CASCADE,
    permissions TEXT
);

CREATE TABLE role_relations
(
    role_id INTEGER REFERENCES roles ON DELETE CASCADE,
    member_id INTEGER REFERENCES members ON DELETE CASCADE,
    PRIMARY KEY (role_id, member_id)
);

CREATE TABLE users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credentials
(
    user_id INTEGER PRIMARY KEY,
    key     TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE
);

CREATE TABLE sessions
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER,
    token      TEXT     NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE
);

CREATE TABLE chats
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    creator_id INTEGER NULL,
    FOREIGN KEY (creator_id) REFERENCES users ON DELETE SET NULL
)

CREATE TABLE members
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER,
    chat_id    INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE CASCADE
)
