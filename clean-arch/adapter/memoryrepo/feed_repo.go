package memoryrepo

import (
	"fmt"
	"sync"

	"rss-aggregator/clean-arch/entity"
)

// InMemoryFeedRepository реализует FeedRepository в памяти
type InMemoryFeedRepository struct {
	feeds  map[int]*entity.Feed
	urls   map[string]int
	mu     sync.RWMutex
	nextID int
}

// NewInMemoryFeedRepository создает новый экземпляр InMemoryFeedRepository
func NewInMemoryFeedRepository() *InMemoryFeedRepository {
	return &InMemoryFeedRepository{
		feeds:  make(map[int]*entity.Feed),
		urls:   make(map[string]int),
		nextID: 1,
	}
}

// Create создает новую RSS-ленту
func (r *InMemoryFeedRepository) Create(feed *entity.Feed) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.urls[feed.URL]; exists {
		return fmt.Errorf("feed with URL %s already exists", feed.URL)
	}

	feed.ID = r.nextID
	r.nextID++
	r.feeds[feed.ID] = feed
	r.urls[feed.URL] = feed.ID

	return nil
}

// GetByURL получает ленту по URL
func (r *InMemoryFeedRepository) GetByURL(url string) (*entity.Feed, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.urls[url]
	if !exists {
		return nil, nil
	}

	feed := *r.feeds[id]
	return &feed, nil
}

// GetAll получает все ленты
func (r *InMemoryFeedRepository) GetAll() ([]*entity.Feed, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	feeds := make([]*entity.Feed, 0, len(r.feeds))
	for _, feed := range r.feeds {
		feedCopy := *feed
		feeds = append(feeds, &feedCopy)
	}

	return feeds, nil
}

// GetByID получает ленту по ID
func (r *InMemoryFeedRepository) GetByID(id int) (*entity.Feed, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	feed, exists := r.feeds[id]
	if !exists {
		return nil, nil
	}

	feedCopy := *feed
	return &feedCopy, nil
}
