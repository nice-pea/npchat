CREATE TABLE chats
(
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    chief_user_id TEXT NOT NULL,
    FOREIGN KEY (chief_user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE members
(
    id      TEXT PRIMARY KEY,
    chat_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE invitations
(
    id              TEXT PRIMARY KEY,
    subject_user_id TEXT NOT NULL,
    user_id         TEXT NOT NULL,
    chat_id         TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT,
    FOREIGN KEY (subject_user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE users
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    nick TEXT NOT NULL
);

CREATE TABLE sessions
(
    id      TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token   TEXT NOT NULL,
    status  INT  NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE authn_passwords
(
    user_id  TEXT PRIMARY KEY,
    login    TEXT NOT NULL,
    password TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE oauth_tokens
(
    access_token  TEXT NOT NULL,
    token_type    TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry        INT  NOT NULL,
    link_id       TEXT NOT NULL,
    FOREIGN KEY (link_id) REFERENCES oauth_links ON DELETE RESTRICT
);

CREATE TABLE oauth_links
(
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL,
    external_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
);