package usecase

import (
	"fmt"

	"rss-aggregator/clean-arch/entity"
)

// AddFeedUseCase представляет use case для добавления RSS-ленты
type AddFeedUseCase struct {
	feedRepo    entity.FeedRepository
	articleRepo entity.ArticleRepository
	parser      entity.RSSParser
}

// NewAddFeedUseCase создает новый экземпляр AddFeedUseCase
func NewAddFeedUseCase(feedRepo entity.FeedRepository, articleRepo entity.ArticleRepository, parser entity.RSSParser) *AddFeedUseCase {
	return &AddFeedUseCase{
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
		parser:      parser,
	}
}

// Execute выполняет добавление RSS-ленты
func (uc *AddFeedUseCase) Execute(url string) (*entity.Feed, error) {
	// Проверяем, не существует ли уже лента с таким URL
	existingFeed, err := uc.feedRepo.GetByURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to check feed existence: %w", err)
	}
	if existingFeed != nil {
		return nil, fmt.Errorf("feed with URL %s already exists", url)
	}

	// Парсим RSS-ленту
	parsedFeed, err := uc.parser.ParseFeed(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	// Создаем ленту
	feed := &entity.Feed{
		URL:         url,
		Title:       parsedFeed.Title,
		Description: parsedFeed.Description,
	}

	if err := uc.feedRepo.Create(feed); err != nil {
		return nil, fmt.Errorf("failed to create feed: %w", err)
	}

	// Сохраняем статьи из ленты
	for _, item := range parsedFeed.Items {
		article := &entity.Article{
			FeedID:          feed.ID,
			Title:           item.Title,
			Content:         item.Content,
			PublicationDate: item.PublicationDate,
			IsRead:          false,
		}
		if err := uc.articleRepo.Create(article); err != nil {
			return nil, fmt.Errorf("failed to create article: %w", err)
		}
	}

	return feed, nil
}
