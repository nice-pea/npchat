CREATE TABLE chats
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE members
(
    id      TEXT PRIMARY KEY,
    chat_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT
);

CREATE TABLE invitations
(
    id      TEXT PRIMARY KEY,
    chat_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT
);