CREATE TABLE chats
(
    id       TEXT PRIMARY KEY,
    name     TEXT NOT NULL,
    chief_id TEXT NOT NULL
);

CREATE TABLE participants
(
    chat_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT
);

CREATE TABLE invitations
(
    id           TEXT PRIMARY KEY,
    chat_id      TEXT NOT NULL,
    subject_id   TEXT NOT NULL,
    recipient_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT
);