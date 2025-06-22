CREATE TABLE sessions
(
    id            TEXT PRIMARY KEY,
    user_id       TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    status        TEXT        NOT NULL,

--     access token

    access_token  TEXT        NOT NULL,
    access_expiry TIMESTAMPTZ NOT NULL,

--     refresh token

    refresh_token  TEXT        NOT NULL,
    refresh_expiry TIMESTAMPTZ NOT NULL
);