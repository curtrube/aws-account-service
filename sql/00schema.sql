CREATE SCHEMA IF NOT EXISTS aws_cloud;

CREATE TYPE aws_cloud.account_joined_method AS ENUM (
    'INVITED',
    'CREATED'
);

CREATE TYPE aws_cloud.account_state AS ENUM (
    'PENDING_ACTIVATION',
    'ACTIVE',
    'SUSPENDED',
    'PENDING_CLOSURE',
    'CLOSED'
);

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS aws_cloud.account(
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    account_id       CHAR(12) NOT NULL UNIQUE,
    account_name     VARCHAR(50) NOT NULL,
    account_arn      TEXT,
    email            CITEXT UNIQUE NOT NULL,
    account_state    aws_cloud.account_state NOT NULL,
    joined_method    aws_cloud.account_joined_method NOT NULL,
    joined_at        TIMESTAMPTZ NOT NULL
);
