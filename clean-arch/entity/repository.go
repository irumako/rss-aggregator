package entity

// FeedRepository определяет интерфейс для работы с RSS-лентами
type FeedRepository interface {
	Create(feed *Feed) error
	GetByURL(url string) (*Feed, error)
	GetAll() ([]*Feed, error)
	GetByID(id int) (*Feed, error)
}

// ArticleRepository определяет интерфейс для работы со статьями
type ArticleRepository interface {
	Create(article *Article) error
	GetByFeedID(feedID int) ([]*Article, error)
	GetAll() ([]*Article, error)
	MarkAsRead(articleID int) error
}
