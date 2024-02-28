CREATE TABLE IF NOT EXISTS Users(
    id serial,
    username VARCHAR(30) NOT NULL,
    passwordHash VARCHAR(255) NOT NULL,

    PRIMARY KEY (id)
);