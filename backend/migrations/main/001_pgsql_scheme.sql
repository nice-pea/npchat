CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credentials
(
    user_id INTEGER PRIMARY KEY,
    key     TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE
);

CREATE TABLE sessions
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    token      TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE
);

CREATE TABLE chats
(
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator_id INTEGER,
    FOREIGN KEY (creator_id) REFERENCES users ON DELETE SET NULL
);

CREATE TABLE members
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    chat_id    INTEGER NOT NULL,
    is_pinned  BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE CASCADE
);

CREATE TABLE roles
(
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    chat_id     INTEGER NOT NULL REFERENCES chats ON DELETE CASCADE,
    permissions SMALLINT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE role_relations
(
    role_id   INTEGER REFERENCES roles ON DELETE CASCADE,
    member_id INTEGER REFERENCES members ON DELETE CASCADE,
    PRIMARY KEY (role_id, member_id)
);

CREATE TABLE messages
(
    id          SERIAL PRIMARY KEY,
    chat_id     INTEGER NOT NULL,
    text        TEXT NOT NULL,
    author_id   INTEGER,
    reply_to_id INTEGER,
    edited_at   TIMESTAMP,
    removed_at  TIMESTAMP,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users ON DELETE SET NULL,
    FOREIGN KEY (reply_to_id) REFERENCES messages ON DELETE SET NULL
);
