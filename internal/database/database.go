package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps database connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection
func New(dsn string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Feed represents a feed in the database
type Feed struct {
	ID          int
	URL         string
	Title       *string
	Description *string
}

// Article represents an article in the database
type Article struct {
	ID             int
	FeedID         int
	Title          string
	Content        *string
	PublicationDate *time.Time
	IsRead         bool
}

// GetFeedByURL retrieves a feed by its URL
func (db *DB) GetFeedByURL(url string) (*Feed, error) {
	var feed Feed
	err := db.conn.QueryRow(
		"SELECT id, url, title, description FROM feeds WHERE url = ?",
		url,
	).Scan(&feed.ID, &feed.URL, &feed.Title, &feed.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

// CreateFeed creates a new feed
func (db *DB) CreateFeed(url string, title *string, description *string) (*Feed, error) {
	result, err := db.conn.Exec(
		"INSERT INTO feeds (url, title, description) VALUES (?, ?, ?)",
		url, title, description,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Feed{
		ID:          int(id),
		URL:         url,
		Title:       title,
		Description: description,
	}, nil
}

// CreateArticle creates a new article
func (db *DB) CreateArticle(feedID int, title string, content *string, publicationDate *time.Time) (*Article, error) {
	result, err := db.conn.Exec(
		"INSERT INTO articles (feed_id, title, content, publication_date) VALUES (?, ?, ?, ?)",
		feedID, title, content, publicationDate,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:             int(id),
		FeedID:         feedID,
		Title:          title,
		Content:        content,
		PublicationDate: publicationDate,
		IsRead:         false,
	}, nil
}

// GetArticlesByFeedID retrieves all articles for a feed
func (db *DB) GetArticlesByFeedID(feedID int) ([]Article, error) {
	rows, err := db.conn.Query(
		"SELECT id, feed_id, title, content, publication_date, is_read FROM articles WHERE feed_id = ? ORDER BY publication_date DESC",
		feedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(
			&article.ID,
			&article.FeedID,
			&article.Title,
			&article.Content,
			&article.PublicationDate,
			&article.IsRead,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, rows.Err()
}

// ArticleExists checks if an article with the same title and feed_id already exists
func (db *DB) ArticleExists(feedID int, title string) (bool, error) {
	var count int
	err := db.conn.QueryRow(
		"SELECT COUNT(*) FROM articles WHERE feed_id = ? AND title = ?",
		feedID, title,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

