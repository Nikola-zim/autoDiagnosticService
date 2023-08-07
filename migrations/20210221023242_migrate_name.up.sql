CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    login VARCHAR UNIQUE,
    password VARCHAR,
    balance INTEGER
);

CREATE TABLE IF NOT EXISTS requests(
   ID SERIAL,
   login VARCHAR(255),
   chat_id INTEGER,
   image_path_name VARCHAR(255),
   detected_path_name VARCHAR(255),
   description TEXT,
   status_code INTEGER NOT NULL
);

CREATE UNIQUE INDEX User_login ON users (login);