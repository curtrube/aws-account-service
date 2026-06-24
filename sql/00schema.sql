CREATE SCHEMA IF NOT EXISTS aws;

CREATE TYPE aws.account_joined_method AS ENUM (
    'INVITED',
    'CREATED'
);

CREATE TYPE aws.account_state AS ENUM (
    'PENDING_ACTIVATION',
    'ACTIVE',
    'SUSPENDED',
    'PENDING_CLOSURE',
    'CLOSED'
);

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS aws.account(
    id char(12) PRIMARY KEY,
    name varchar(50) NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    joined_method aws.account_joined_method NOT NULL,
    joined_timestamp timestamptz NOT NULL,
    state aws.account_state NOT NULL
);
