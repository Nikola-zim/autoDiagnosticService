CREATE TABLE IF NOT EXISTS requests(
    ID SERIAL PRIMARY KEY,
    chatID INTEGER NOT NULL,
    image_path_name VARCHAR(255),
    detected_path_name VARCHAR(255),
    description TEXT,
    status_code INTEGER NOT NULL
);
