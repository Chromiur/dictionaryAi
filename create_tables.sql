DROP DATABASE IF EXISTS dictionary_db;
CREATE DATABASE dictionary_db;

DROP TABLE IF EXISTS WordsList;
DROP TABLE IF EXISTS Words;

\c dictionary_db;

CREATE TABLE Words (
    ID SERIAL,
    Word VARCHAR(100)
);

CREATE TABLE WordsList (
    ID Int NOT NULL,
    UserID Int NOT NULL,
    Word VARCHAR(100),
    Description VARCHAR(200)
);

INSERT INTO WordsList (ID, UserID, Word, Description) VALUES (1, 1, 'Hello', 'Привет')