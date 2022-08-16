SELECT 'CREATE DATABASE chat'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'chat')\gexec

\c chat;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_history(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id),
    message VARCHAR,
    created_at TIMESTAMP NOT NULL,
    room VARCHAR NOT NULL
);