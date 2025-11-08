package rss

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

// Parser handles RSS feed parsing
type Parser struct {
	fp *gofeed.Parser
}

// NewParser creates a new RSS parser
func NewParser() *Parser {
	return &Parser{
		fp: gofeed.NewParser(),
	}
}

// FeedInfo represents parsed feed information
type FeedInfo struct {
	Title       string
	Description string
	Items       []Item
}

// Item represents a parsed RSS item
type Item struct {
	Title         string
	Content       string
	PublicationDate *time.Time
}

// ParseFeed parses an RSS feed from a URL
func (p *Parser) ParseFeed(url string) (*FeedInfo, error) {
	feed, err := p.fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	feedInfo := &FeedInfo{
		Title:       feed.Title,
		Description: feed.Description,
		Items:       make([]Item, 0, len(feed.Items)),
	}

	for _, item := range feed.Items {
		var pubDate *time.Time
		if item.PublishedParsed != nil {
			pubDate = item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			pubDate = item.UpdatedParsed
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		feedInfo.Items = append(feedInfo.Items, Item{
			Title:         item.Title,
			Content:       content,
			PublicationDate: pubDate,
		})
	}

	return feedInfo, nil
}

