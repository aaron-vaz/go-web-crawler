package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/aaron-vaz/go-web-crawler/pkg/crawler"
	"github.com/aaron-vaz/go-web-crawler/pkg/httpservice"
	"github.com/aaron-vaz/go-web-crawler/pkg/linksservice"
	"github.com/aaron-vaz/go-web-crawler/pkg/worker"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type App struct {
	td *worker.TaskDispactcher
}

func (app *App) Run() {
	results := app.td.Start()

	fmt.Println()
	fmt.Println("###### Crawler Results ######")
	fmt.Println()

	for _, result := range results {
		fmt.Println(result.FormatResult())
	}
}

func main() {
	var startURL string
	var maxWorkers int

	flag.StringVar(&startURL, "startURL", "", "Crawler starting url, required")
	flag.IntVar(&maxWorkers, "maxWorkers", 10, "Max workers to use to perform the crawling")

	flag.Parse()

	if startURL == "" {
		fmt.Println("No -startURL supplied, exiting...")
		os.Exit(1)
	}

	startURLParsed, err := url.Parse(startURL)
	if err != nil {
		log.Printf("Invalid -startURL provided %s, error = %s", startURL, err.Error())
		os.Exit(1)
	}

	hs := httpservice.NewHttpService(http.DefaultClient)
	ls := linksservice.NewLinksParser()
	wc := crawler.NewWebCrawler(hs, ls)
	td := worker.NewTaskDispactcher(startURLParsed, maxWorkers, wc)

	app := &App{td: td}

	app.Run()
}
