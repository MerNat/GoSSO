#!/bin/sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE sso;
    CREATE USER sso;
    create table users (
        id         serial primary key,
        uuid       varchar(64) not null unique,
        firstname  varchar(64) not null,
        email      varchar(40) not null unique,
        password   varchar(255) not null,
        created_at timestamp not null   
    );
    GRANT ALL ON DATABASE sso TO "sso";
    ALTER USER sso PASSWORD 'password';
    ALTER USER sso CREATEDB;
EOSQL