package app

import (
	"rss-aggregator/clean-arch/adapter"
	"rss-aggregator/clean-arch/adapter/cli"
	"rss-aggregator/clean-arch/adapter/memoryrepo"
	"rss-aggregator/clean-arch/usecase"
)

// App представляет приложение RSS-агрегатора
type App struct {
	cli *cli.CLI
}

// NewApp создает новое приложение
func NewApp() *App {
	// Инициализация репозиториев
	feedRepo := memoryrepo.NewInMemoryFeedRepository()
	articleRepo := memoryrepo.NewInMemoryArticleRepository()

	// Инициализация адаптера RSS-парсера
	rssParser := adapter.NewRSSParserAdapter()

	// Инициализация use cases
	addFeedUseCase := usecase.NewAddFeedUseCase(feedRepo, articleRepo, rssParser)
	listFeedsUseCase := usecase.NewListFeedsUseCase(feedRepo)
	fetchArticlesUseCase := usecase.NewFetchArticlesUseCase(feedRepo, articleRepo, rssParser)
	listArticlesUseCase := usecase.NewListArticlesUseCase(articleRepo)

	// Инициализация CLI
	cliInstance := cli.NewCLI(
		addFeedUseCase,
		listFeedsUseCase,
		fetchArticlesUseCase,
		listArticlesUseCase,
	)

	return &App{
		cli: cliInstance,
	}
}

// Run запускает приложение
func (a *App) Run() {
	a.cli.Run()
}
