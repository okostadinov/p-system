CREATE DATABASE IF NOT EXISTS p_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE p_system;

DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS medications;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

CREATE TABLE medications (name VARCHAR(30) NOT NULL PRIMARY KEY);

CREATE TABLE patients (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    ucn VARCHAR(10) NOT NULL,
    first_name VARCHAR(20) NOT NULL,
    last_name VARCHAR(20) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    height VARCHAR(3) NOT NULL,
    weight VARCHAR(3) NOT NULL,
    medication VARCHAR(30) NOT NULL,
    note TEXT NOT NULL,
    approved BOOLEAN NOT NULL DEFAULT 0,
    first_continuation BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (medication) REFERENCES medications(name)
);

CREATE TABLE sessions (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    session_data LONGBLOB,
    created_on TIMESTAMP DEFAULT NOW(),
    modified_on TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_on TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE
    users
ADD
    CONSTRAINT users_uc_email UNIQUE (email);

DROP USER IF EXISTS 'p_system_admin'@'localhost';

CREATE USER 'p_system_admin' @'localhost';

GRANT
SELECT
,
INSERT
,
UPDATE
,
    DELETE ON p_system.* TO 'p_system_admin' @'localhost';

ALTER USER 'p_system_admin' @'localhost' IDENTIFIED BY 'p_system_admin';