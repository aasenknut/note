BEGIN TRANSACTION;

CREATE TABLE note (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    created DATETIME NOT NULL DEFAULT (DATETIME(CURRENT_TIMESTAMP, 'localtime'))
);

CREATE INDEX created_idx ON note (created);

INSERT INTO note (title, text) VALUES ("First", "This is the very first note written into this app");
INSERT INTO note (title, text) VALUES ("2nd", "Second note written.");
INSERT INTO note (title, text) VALUES ("Short", "Shortest.");
INSERT INTO note (title, text) VALUES ("Today", "Today I did alot of work. It was fun. Looking forward to tomorrow. Tomorrow will be fun too.");

CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    created DATETIME NOT NULL DEFAULT (DATETIME(CURRENT_TIMESTAMP, 'localtime'))
);

INSERT INTO user (username) VALUES ("user");

CREATE TABLE auth (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id TEXT NOT NULL,
    password TEXT NOT NULL,
    created DATETIME NOT NULL DEFAULT (DATETIME(CURRENT_TIMESTAMP, 'localtime'))
);

COMMIT;
