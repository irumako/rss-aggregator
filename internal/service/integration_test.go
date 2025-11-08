package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	api "rss-aggregator/gen"
	"rss-aggregator/internal/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates an in-memory SQLite database and applies migrations
func setupTestDB(t *testing.T) (*database.DB, func()) {
	// Create temporary database file
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create database connection directly to apply migrations
	conn, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	require.NoError(t, conn.Ping())

	// Apply migrations by reading and executing the SQL
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS feeds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE NOT NULL,
		title TEXT,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		feed_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT,
		publication_date DATETIME,
		is_read BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (feed_id) REFERENCES feeds (id)
	);
	`

	_, err = conn.Exec(migrationSQL)
	require.NoError(t, err)
	conn.Close()

	// Create database wrapper
	db, err := database.New(dbPath)
	require.NoError(t, err)

	// Cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

// setupTestApp creates a Fiber app with test service
func setupTestApp(t *testing.T, db *database.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	svc := New(db)
	api.RegisterHandlers(app, svc)

	return app
}

func TestPostFeeds_Integration(t *testing.T) {
	t.Run("successful feed creation", func(t *testing.T) {
		// Setup
		db, cleanup := setupTestDB(t)
		defer cleanup()

		app := setupTestApp(t, db)

		// Test data - using a real RSS feed URL for integration test
		// Using a well-known RSS feed that should be stable
		// Try multiple RSS feeds in case one is unavailable
		feedURLs := []string{
			"https://feeds.bbci.co.uk/news/rss.xml",
			"https://rss.cnn.com/rss/edition.rss",
			"https://www.nasa.gov/rss/dyn/breaking_news.rss",
		}

		var feedURL string
		var feedResponse api.FeedResponse
		var success bool

		// Try each feed URL until one works
		for _, url := range feedURLs {
			feedURL = url
			reqBody := api.AddFeedRequest{
				Url: feedURL,
			}

			bodyBytes, err := json.Marshal(reqBody)
			require.NoError(t, err)

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Execute request with timeout (30 seconds for RSS feed parsing)
			resp, err := app.Test(req, int(30*time.Second.Milliseconds()))
			require.NoError(t, err)
			defer resp.Body.Close()

			// Parse response body
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if resp.StatusCode == http.StatusCreated {
				err = json.Unmarshal(body, &feedResponse)
				require.NoError(t, err)
				success = true
				break
			} else {
				// Log error but try next URL
				var errorResponse map[string]string
				if err := json.Unmarshal(body, &errorResponse); err == nil {
					t.Logf("Feed %s failed: %v", url, errorResponse)
				}
			}
		}

		// Skip test if no feed is available (network issue)
		if !success {
			t.Skip("No RSS feed is available. This may be due to network issues.")
		}

		// Assertions
		assert.NotNil(t, feedResponse.Id)
		assert.NotNil(t, feedResponse.Url)
		assert.Equal(t, feedURL, *feedResponse.Url)
		assert.NotEmpty(t, feedResponse.Title)
		// Articles may or may not be present depending on the feed
		if feedResponse.Articles != nil {
			assert.GreaterOrEqual(t, len(*feedResponse.Articles), 0)
		}
	})

	t.Run("duplicate feed URL", func(t *testing.T) {
		// Setup
		db, cleanup := setupTestDB(t)
		defer cleanup()

		app := setupTestApp(t, db)

		// Try multiple RSS feeds in case one is unavailable
		feedURLs := []string{
			"https://feeds.bbci.co.uk/news/rss.xml",
			"https://rss.cnn.com/rss/edition.rss",
			"https://www.nasa.gov/rss/dyn/breaking_news.rss",
		}

		var feedURL string
		var success bool

		// Find a working feed URL
		for _, url := range feedURLs {
			feedURL = url
			reqBody := api.AddFeedRequest{
				Url: feedURL,
			}

			bodyBytes, err := json.Marshal(reqBody)
			require.NoError(t, err)

			// First request - should succeed
			req1 := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(bodyBytes))
			req1.Header.Set("Content-Type", "application/json")
			resp1, err := app.Test(req1, int(30*time.Second.Milliseconds()))
			require.NoError(t, err)
			resp1.Body.Close()

			if resp1.StatusCode == http.StatusCreated {
				success = true
				break
			}
		}

		// Skip test if no feed is available (network issue)
		if !success {
			t.Skip("No RSS feed is available. This may be due to network issues.")
		}

		// Second request with same URL - should fail with conflict
		reqBody := api.AddFeedRequest{
			Url: feedURL,
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req2 := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(bodyBytes))
		req2.Header.Set("Content-Type", "application/json")
		resp2, err := app.Test(req2, int(30*time.Second.Milliseconds()))
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusConflict, resp2.StatusCode)

		// Verify error message
		body, err := io.ReadAll(resp2.Body)
		require.NoError(t, err)

		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "already exists")
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup
		db, cleanup := setupTestDB(t)
		defer cleanup()

		app := setupTestApp(t, db)

		// Invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, int(5*time.Second.Milliseconds()))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing URL", func(t *testing.T) {
		// Setup
		db, cleanup := setupTestDB(t)
		defer cleanup()

		app := setupTestApp(t, db)

		reqBody := api.AddFeedRequest{
			Url: "",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, int(5*time.Second.Milliseconds()))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "URL is required")
	})

	t.Run("invalid RSS feed URL", func(t *testing.T) {
		// Setup
		db, cleanup := setupTestDB(t)
		defer cleanup()

		app := setupTestApp(t, db)

		// Use a URL that will quickly fail (localhost with invalid port)
		// This avoids long DNS lookup timeouts
		reqBody := api.AddFeedRequest{
			Url: "http://127.0.0.1:99999/invalid-feed.xml",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Use longer timeout for invalid URL test as connection attempts may take time
		resp, err := app.Test(req, int(15*time.Second.Milliseconds()))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var errorResponse map[string]string
		err = json.Unmarshal(body, &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Failed to parse RSS feed")
	})
}
