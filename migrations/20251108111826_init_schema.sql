-- +goose Up
-- +goose StatementBegin
-- Таблица пользователей
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL
);

-- Таблица категорий
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Таблица фидов
CREATE TABLE feeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT UNIQUE NOT NULL,
    title TEXT,
    description TEXT
);

-- Таблица статей
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    publication_date DATETIME,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- Связь пользователь-фид
CREATE TABLE user_feeds (
    user_id INTEGER,
    feed_id INTEGER,
    PRIMARY KEY (user_id, feed_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (feed_id) REFERENCES feeds (id)
);

-- Связь фид-категория
CREATE TABLE feed_categories (
    feed_id INTEGER,
    category_id INTEGER,
    PRIMARY KEY (feed_id, category_id),
    FOREIGN KEY (feed_id) REFERENCES feeds (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feed_categories;
DROP TABLE IF EXISTS user_feeds;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS feeds;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
