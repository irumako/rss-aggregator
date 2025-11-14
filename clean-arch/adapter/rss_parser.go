package adapter

import (
	"rss-aggregator/clean-arch/entity"
	"rss-aggregator/internal/rss"
)

// RSSParserAdapter адаптирует internal/rss.Parser к entity.RSSParser
type RSSParserAdapter struct {
	parser *rss.Parser
}

// NewRSSParserAdapter создает новый экземпляр RSSParserAdapter
func NewRSSParserAdapter() *RSSParserAdapter {
	return &RSSParserAdapter{
		parser: rss.NewParser(),
	}
}

// ParseFeed парсит RSS-ленту
func (a *RSSParserAdapter) ParseFeed(url string) (*entity.ParsedFeed, error) {
	feedInfo, err := a.parser.ParseFeed(url)
	if err != nil {
		return nil, err
	}

	parsedFeed := &entity.ParsedFeed{
		Title:       feedInfo.Title,
		Description: feedInfo.Description,
		Items:       make([]entity.ParsedItem, 0, len(feedInfo.Items)),
	}

	for _, item := range feedInfo.Items {
		parsedFeed.Items = append(parsedFeed.Items, entity.ParsedItem{
			Title:           item.Title,
			Content:         item.Content,
			PublicationDate: item.PublicationDate,
		})
	}

	return parsedFeed, nil
}
