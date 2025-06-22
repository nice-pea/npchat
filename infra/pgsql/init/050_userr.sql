CREATE TABLE users
(
    id       TEXT PRIMARY KEY,
    name     TEXT NOT NULL,
    nick     TEXT NOT NULL,
    login    TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE oauth_users
(
    id            TEXT PRIMARY KEY,
    user_id       TEXT        NOT NULL,
    provider      TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    picture       TEXT        NOT NULL,

--     token

    access_token  TEXT        NOT NULL,
    token_type    TEXT        NOT NULL,
    refresh_token TEXT        NOT NULL,
    expiry        TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
);