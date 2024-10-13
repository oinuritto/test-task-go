-- DROP TABLE refresh_tokens;
-- DROP TABLE users;

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY,
    email    VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    user_id    UUID REFERENCES users (id) PRIMARY KEY ,
    token_hash TEXT NOT NULL,
    ip_address VARCHAR(45),
    created_at TIMESTAMP        DEFAULT NOW(),
    expires_at TIMESTAMP
);
