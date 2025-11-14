package entity

import "time"

// RSSParser определяет интерфейс для парсинга RSS-лент
type RSSParser interface {
	ParseFeed(url string) (*ParsedFeed, error)
}

// ParsedFeed представляет распарсенную RSS-ленту
type ParsedFeed struct {
	Title       string
	Description string
	Items       []ParsedItem
}

// ParsedItem представляет распарсенную статью из RSS
type ParsedItem struct {
	Title           string
	Content         string
	PublicationDate *time.Time
}
