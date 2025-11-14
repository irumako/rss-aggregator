package memoryrepo

import (
	"fmt"
	"sync"

	"rss-aggregator/clean-arch/entity"
)

// InMemoryArticleRepository реализует ArticleRepository в памяти
type InMemoryArticleRepository struct {
	articles map[int]*entity.Article
	feedArticles map[int][]int
	mu       sync.RWMutex
	nextID   int
}

// NewInMemoryArticleRepository создает новый экземпляр InMemoryArticleRepository
func NewInMemoryArticleRepository() *InMemoryArticleRepository {
	return &InMemoryArticleRepository{
		articles:    make(map[int]*entity.Article),
		feedArticles: make(map[int][]int),
		nextID:      1,
	}
}

// Create создает новую статью
func (r *InMemoryArticleRepository) Create(article *entity.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	article.ID = r.nextID
	r.nextID++
	r.articles[article.ID] = article
	r.feedArticles[article.FeedID] = append(r.feedArticles[article.FeedID], article.ID)

	return nil
}

// GetByFeedID получает все статьи для ленты
func (r *InMemoryArticleRepository) GetByFeedID(feedID int) ([]*entity.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	articleIDs, exists := r.feedArticles[feedID]
	if !exists {
		return []*entity.Article{}, nil
	}

	articles := make([]*entity.Article, 0, len(articleIDs))
	for _, id := range articleIDs {
		article := *r.articles[id]
		articles = append(articles, &article)
	}

	return articles, nil
}

// GetAll получает все статьи
func (r *InMemoryArticleRepository) GetAll() ([]*entity.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	articles := make([]*entity.Article, 0, len(r.articles))
	for _, article := range r.articles {
		articleCopy := *article
		articles = append(articles, &articleCopy)
	}

	return articles, nil
}

// MarkAsRead помечает статью как прочитанную
func (r *InMemoryArticleRepository) MarkAsRead(articleID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	article, exists := r.articles[articleID]
	if !exists {
		return fmt.Errorf("article with ID %d not found", articleID)
	}

	article.IsRead = true
	return nil
}

