package usecase

import "rss-aggregator/clean-arch/entity"

// ListFeedsUseCase представляет use case для получения списка RSS-лент
type ListFeedsUseCase struct {
	feedRepo entity.FeedRepository
}

// NewListFeedsUseCase создает новый экземпляр ListFeedsUseCase
func NewListFeedsUseCase(feedRepo entity.FeedRepository) *ListFeedsUseCase {
	return &ListFeedsUseCase{
		feedRepo: feedRepo,
	}
}

// Execute выполняет получение списка всех RSS-лент
func (uc *ListFeedsUseCase) Execute() ([]*entity.Feed, error) {
	return uc.feedRepo.GetAll()
}

