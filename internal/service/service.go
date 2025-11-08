package service

import (
	"fmt"

	api "rss-aggregator/gen"
	"rss-aggregator/internal/database"
	"rss-aggregator/internal/rss"

	"github.com/gofiber/fiber/v2"
)

// Service implements the ServerInterface
type Service struct {
	db     *database.DB
	parser *rss.Parser
}

// New creates a new service instance
func New(db *database.DB) *Service {
	return &Service{
		db:     db,
		parser: rss.NewParser(),
	}
}

// PostFeeds handles POST /feeds request
func (s *Service) PostFeeds(c *fiber.Ctx) error {
	var req api.AddFeedRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Url == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL is required",
		})
	}

	// Check if feed already exists
	existingFeed, err := s.db.GetFeedByURL(req.Url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check feed existence",
		})
	}

	if existingFeed != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Feed with this URL already exists",
		})
	}

	// Parse RSS feed
	feedInfo, err := s.parser.ParseFeed(req.Url)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to parse RSS feed: %v", err),
		})
	}

	// Create feed in database
	title := feedInfo.Title
	description := feedInfo.Description
	feed, err := s.db.CreateFeed(req.Url, &title, &description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create feed",
		})
	}

	// Save articles from RSS feed
	for _, item := range feedInfo.Items {
		// Check if article already exists
		exists, err := s.db.ArticleExists(feed.ID, item.Title)
		if err != nil {
			continue // Skip on error
		}
		if exists {
			continue // Skip duplicate
		}

		// Create article
		_, err = s.db.CreateArticle(
			feed.ID,
			item.Title,
			&item.Content,
			item.PublicationDate,
		)
		if err != nil {
			continue // Skip on error
		}
	}

	// Get all articles for the feed
	allArticles, err := s.db.GetArticlesByFeedID(feed.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve articles",
		})
	}

	// Convert to API model
	articles := make([]api.Article, 0, len(allArticles))
	for _, article := range allArticles {
		apiArticle := api.Article{
			Id:              &article.ID,
			Title:           &article.Title,
			Content:         article.Content,
			PublicationDate: article.PublicationDate,
			IsRead:          &article.IsRead,
		}
		articles = append(articles, apiArticle)
	}

	// Build response
	response := api.FeedResponse{
		Id:          &feed.ID,
		Url:         &feed.URL,
		Title:       feed.Title,
		Description: feed.Description,
		Articles:    &articles,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Ensure Service implements ServerInterface
var _ api.ServerInterface = (*Service)(nil)

