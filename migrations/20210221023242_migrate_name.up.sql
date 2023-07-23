CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    login VARCHAR UNIQUE,
    password VARCHAR
);

CREATE TABLE IF NOT EXISTS requests(
   ID SERIAL,
   user_id INTEGER,
   FOREIGN KEY (user_id) REFERENCES users(user_id),
   chat_id INTEGER,
   image_path_name VARCHAR(255),
   detected_path_name VARCHAR(255),
   description TEXT,
   status_code INTEGER NOT NULL
);

CREATE UNIQUE INDEX User_login ON users (login);