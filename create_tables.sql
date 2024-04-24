DROP
DATABASE IF EXISTS dictionary_db;
CREATE
DATABASE dictionary_db;

DROP TABLE IF EXISTS WordsList;
DROP TABLE IF EXISTS Words;

\c
dictionary_db;

CREATE TABLE Words
(
    ID   SERIAL,
    Word VARCHAR(100)
);

CREATE TABLE WordsList
(
    ID          Int NOT NULL,
    UserID      Int NOT NULL,
    Word        VARCHAR(100),
    Description VARCHAR(200)
);

/*
Users
*/
CREATE TABLE users
(
    id        SERIAL PRIMARY KEY,
    username  VARCHAR NOT NULL,
    password  VARCHAR NOT NULL,
    lastlogin INT,
    admin     INT,
    active    INT,
    teacher   INT,
    student   INT
);

INSERT INTO users (username, password, lastlogin, admin, active, teacher, student)
VALUES ('admin', 'admin', 1620922454, 1, 1, 0, 0);