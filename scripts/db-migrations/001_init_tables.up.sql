CREATE TYPE role AS ENUM ('admin', 'usermanager', 'user');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name TEXT,
    role role NOT NULL,
    created_at TIMESTAMP DEFAULT (now()),
    updated_at TIMESTAMP DEFAULT (now())
);

CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    ts TIMESTAMP NOT NULL,
    loc TEXT NOT NULL,
    distance INT NOT NULL,
    created_at TIMESTAMP DEFAULT (now()),
    updated_at TIMESTAMP DEFAULT (now())
);
ALTER TABLE activities ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
CREATE INDEX activities_ts_idx ON activities(ts);
