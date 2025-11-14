package entity

import "time"

// Feed представляет RSS-ленту
type Feed struct {
	ID          int
	URL         string
	Title       string
	Description string
}

// Article представляет статью из RSS-ленты
type Article struct {
	ID             int
	FeedID         int
	Title          string
	Content        string
	PublicationDate *time.Time
	IsRead         bool
}

