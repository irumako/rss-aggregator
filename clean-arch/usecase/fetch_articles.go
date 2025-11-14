package usecase

import (
	"fmt"

	"rss-aggregator/clean-arch/entity"
)

// FetchArticlesUseCase представляет use case для обновления статей из RSS-ленты
type FetchArticlesUseCase struct {
	feedRepo    entity.FeedRepository
	articleRepo entity.ArticleRepository
	parser      entity.RSSParser
}

// NewFetchArticlesUseCase создает новый экземпляр FetchArticlesUseCase
func NewFetchArticlesUseCase(feedRepo entity.FeedRepository, articleRepo entity.ArticleRepository, parser entity.RSSParser) *FetchArticlesUseCase {
	return &FetchArticlesUseCase{
		feedRepo:    feedRepo,
		articleRepo: articleRepo,
		parser:      parser,
	}
}

// Execute выполняет обновление статей из RSS-ленты
func (uc *FetchArticlesUseCase) Execute(feedID int) error {
	// Получаем ленту
	feed, err := uc.feedRepo.GetByID(feedID)
	if err != nil {
		return fmt.Errorf("failed to get feed: %w", err)
	}
	if feed == nil {
		return fmt.Errorf("feed with ID %d not found", feedID)
	}

	// Парсим RSS-ленту
	parsedFeed, err := uc.parser.ParseFeed(feed.URL)
	if err != nil {
		return fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	// Получаем существующие статьи
	existingArticles, err := uc.articleRepo.GetByFeedID(feedID)
	if err != nil {
		return fmt.Errorf("failed to get existing articles: %w", err)
	}

	// Создаем map для быстрой проверки существующих статей
	existingTitles := make(map[string]bool)
	for _, article := range existingArticles {
		existingTitles[article.Title] = true
	}

	// Добавляем новые статьи
	for _, item := range parsedFeed.Items {
		if !existingTitles[item.Title] {
			article := &entity.Article{
				FeedID:         feedID,
				Title:          item.Title,
				Content:        item.Content,
				PublicationDate: item.PublicationDate,
				IsRead:         false,
			}
			if err := uc.articleRepo.Create(article); err != nil {
				return fmt.Errorf("failed to create article: %w", err)
			}
		}
	}

	return nil
}

