package usecase

import "rss-aggregator/clean-arch/entity"

// ListArticlesUseCase представляет use case для получения списка статей
type ListArticlesUseCase struct {
	articleRepo entity.ArticleRepository
}

// NewListArticlesUseCase создает новый экземпляр ListArticlesUseCase
func NewListArticlesUseCase(articleRepo entity.ArticleRepository) *ListArticlesUseCase {
	return &ListArticlesUseCase{
		articleRepo: articleRepo,
	}
}

// Execute выполняет получение списка статей
// Если feedID > 0, возвращает статьи только для этой ленты
func (uc *ListArticlesUseCase) Execute(feedID int) ([]*entity.Article, error) {
	if feedID > 0 {
		return uc.articleRepo.GetByFeedID(feedID)
	}
	return uc.articleRepo.GetAll()
}

