package main

import (
	"rss-aggregator/clean-arch/app"
)

func main() {
	application := app.NewApp()
	application.Run()
}

