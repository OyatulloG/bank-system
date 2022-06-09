DROP DATABASE IF EXISTS bank;
CREATE DATABASE bank;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    userid BIGSERIAL PRIMARY KEY NOT NULL,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    regtime DATE DEFAULT CURRENT_DATE,
    pin VARCHAR(4) NOT NULL,
    balance NUMERIC(19,2) DEFAULT 0,
    status BOOLEAN DEFAULT 't',
    dateofbirth DATE NOT NULL,
    gender VARCHAR(6) NOT NULL,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

INSERT INTO users (firstname, lastname, pin, dateofbirth, gender, username, password)
VALUES ('John', 'Smith', '1234', '1995-10-31', 'male', 'JohnS', 'doc@1995');

INSERT INTO users (firstname, lastname, pin, dateofbirth, gender, username, password)
VALUES ('Anna', 'Smith', '4321', '1996-02-09', 'female', 'AnnaS', 'doc@1996');



DROP TABLE IF EXISTS admins;
CREATE TABLE admins (
    adminid BIGSERIAL PRIMARY KEY NOT NULL,
    password VARCHAR(50) NOT NULL
);

INSERT INTO admins (password)
VALUES ('root@123');

INSERT INTO admins (password)
VALUES ('abc@123');



DROP TABLE IF EXISTS messages;
CREATE TABLE messages (
    messageid BIGSERIAL PRIMARY KEY NOT NULL,
    userid BIGSERIAL NOT NULL,
    adminid BIGSERIAL NOT NULL,
    ischecked BOOLEAN DEFAULT 'f',
    date DATE DEFAULT CURRENT_DATE,
    message TEXT
);

INSERT INTO messages (userid, adminid, message)
VALUES (1, 1, 'problem 1');
INSERT INTO messages (userid, adminid, message)
VALUES (2, 1, 'problem 2');



DROP TABLE IF EXISTS transfers;
CREATE TABLE transfers (
    transferid BIGSERIAL PRIMARY KEY NOT NULL,
    senderid BIGSERIAL NOT NULL,
    recieverid BIGSERIAL NOT NULL,
    amount NUMERIC(19, 2) DEFAULT 0,
    date DATE DEFAULT CURRENT_DATE
);

INSERT INTO transfers (senderid, recieverid, amount)
VALUES (1, 1, 10000);
INSERT INTO transfers (senderid, recieverid, amount)
VALUES (1, 2, 5000.50);
INSERT INTO transfers (senderid, recieverid, amount)
VALUES (2, 1, 5000);