CREATE TABLE chats
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    chief_user_id TEXT NOT NULL
);

CREATE TABLE members
(
    id      TEXT PRIMARY KEY,
    chat_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE invitations
(
    id      TEXT PRIMARY KEY,
    subject_user_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT
    FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT
    FOREIGN KEY (subject_user_id) REFERENCES users ON DELETE RESTRICT
);

CREATE TABLE users
(
    id      TEXT PRIMARY KEY
);