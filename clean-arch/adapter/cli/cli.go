package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"rss-aggregator/clean-arch/usecase"
)

// CLI представляет интерфейс командной строки
type CLI struct {
	addFeedUseCase      *usecase.AddFeedUseCase
	listFeedsUseCase    *usecase.ListFeedsUseCase
	fetchArticlesUseCase *usecase.FetchArticlesUseCase
	listArticlesUseCase  *usecase.ListArticlesUseCase
	scanner             *bufio.Scanner
}

// NewCLI создает новый экземпляр CLI
func NewCLI(
	addFeedUseCase *usecase.AddFeedUseCase,
	listFeedsUseCase *usecase.ListFeedsUseCase,
	fetchArticlesUseCase *usecase.FetchArticlesUseCase,
	listArticlesUseCase *usecase.ListArticlesUseCase,
) *CLI {
	return &CLI{
		addFeedUseCase:      addFeedUseCase,
		listFeedsUseCase:    listFeedsUseCase,
		fetchArticlesUseCase: fetchArticlesUseCase,
		listArticlesUseCase:  listArticlesUseCase,
		scanner:             bufio.NewScanner(os.Stdin),
	}
}

// Run запускает CLI
func (c *CLI) Run() {
	fmt.Println("=== RSS Aggregator ===")
	fmt.Println("Доступные команды:")
	fmt.Println("  add <url>          - Добавить RSS-ленту")
	fmt.Println("  list-feeds         - Показать все RSS-ленты")
	fmt.Println("  fetch <feed-id>    - Обновить статьи из ленты")
	fmt.Println("  articles [feed-id] - Показать статьи (опционально для конкретной ленты)")
	fmt.Println("  help               - Показать эту справку")
	fmt.Println("  exit               - Выход")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !c.scanner.Scan() {
			break
		}

		line := strings.TrimSpace(c.scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Ошибка: укажите URL RSS-ленты")
				continue
			}
			c.handleAddFeed(parts[1])

		case "list-feeds":
			c.handleListFeeds()

		case "fetch":
			if len(parts) < 2 {
				fmt.Println("Ошибка: укажите ID ленты")
				continue
			}
			feedID, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Ошибка: неверный ID ленты: %v\n", err)
				continue
			}
			c.handleFetchArticles(feedID)

		case "articles":
			feedID := 0
			if len(parts) >= 2 {
				var err error
				feedID, err = strconv.Atoi(parts[1])
				if err != nil {
					fmt.Printf("Ошибка: неверный ID ленты: %v\n", err)
					continue
				}
			}
			c.handleListArticles(feedID)

		case "help":
			c.printHelp()

		case "exit":
			fmt.Println("До свидания!")
			return

		default:
			fmt.Printf("Неизвестная команда: %s. Введите 'help' для справки.\n", command)
		}
	}
}

func (c *CLI) handleAddFeed(url string) {
	feed, err := c.addFeedUseCase.Execute(url)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("✓ RSS-лента добавлена:\n")
	fmt.Printf("  ID: %d\n", feed.ID)
	fmt.Printf("  URL: %s\n", feed.URL)
	fmt.Printf("  Название: %s\n", feed.Title)
	if feed.Description != "" {
		fmt.Printf("  Описание: %s\n", feed.Description)
	}
	fmt.Println()
}

func (c *CLI) handleListFeeds() {
	feeds, err := c.listFeedsUseCase.Execute()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	if len(feeds) == 0 {
		fmt.Println("Нет добавленных RSS-лент")
		fmt.Println()
		return
	}

	fmt.Println("RSS-ленты:")
	for _, feed := range feeds {
		fmt.Printf("  [%d] %s\n", feed.ID, feed.Title)
		fmt.Printf("      URL: %s\n", feed.URL)
		if feed.Description != "" {
			fmt.Printf("      Описание: %s\n", feed.Description)
		}
		fmt.Println()
	}
}

func (c *CLI) handleFetchArticles(feedID int) {
	err := c.fetchArticlesUseCase.Execute(feedID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("✓ Статьи из ленты %d обновлены\n", feedID)
	fmt.Println()
}

func (c *CLI) handleListArticles(feedID int) {
	articles, err := c.listArticlesUseCase.Execute(feedID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	if len(articles) == 0 {
		if feedID > 0 {
			fmt.Printf("Нет статей для ленты %d\n", feedID)
		} else {
			fmt.Println("Нет статей")
		}
		fmt.Println()
		return
	}

	if feedID > 0 {
		fmt.Printf("Статьи из ленты %d:\n", feedID)
	} else {
		fmt.Println("Все статьи:")
	}

	for _, article := range articles {
		readStatus := " "
		if article.IsRead {
			readStatus = "✓"
		}
		fmt.Printf("  [%s] [%d] %s\n", readStatus, article.ID, article.Title)
		if article.PublicationDate != nil {
			fmt.Printf("      Дата: %s\n", article.PublicationDate.Format("2006-01-02 15:04:05"))
		}
		if article.Content != "" && len(article.Content) > 100 {
			fmt.Printf("      %s...\n", article.Content[:100])
		} else if article.Content != "" {
			fmt.Printf("      %s\n", article.Content)
		}
		fmt.Println()
	}
}

func (c *CLI) printHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  add <url>          - Добавить RSS-ленту")
	fmt.Println("  list-feeds         - Показать все RSS-ленты")
	fmt.Println("  fetch <feed-id>    - Обновить статьи из ленты")
	fmt.Println("  articles [feed-id] - Показать статьи (опционально для конкретной ленты)")
	fmt.Println("  help               - Показать эту справку")
	fmt.Println("  exit               - Выход")
	fmt.Println()
}

