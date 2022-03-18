package worker

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aaron-vaz/go-web-crawler/pkg/crawler"
	"github.com/aaron-vaz/go-web-crawler/pkg/httpservice"
)

const (
	backOffTimeout = 1 * time.Minute
)

type TaskResult struct {
	url   *url.URL
	links []*url.URL
}

func (tr TaskResult) FormatResult() string {
	result := fmt.Sprintf("=========== Links for %s =====>\n", tr.url)

	for _, link := range tr.links {
		result += link.String() + "\n"
	}

	if len(tr.links) == 0 {
		result += fmt.Sprintln("No links found on page")
	}

	return result
}

type Task struct {
	url     *url.URL
	links   chan<- *url.URL
	results chan<- TaskResult

	wc crawler.WebCrawler
}

func (t Task) Run() {
	links, err := t.wc.Crawl(t.url)

	if errors.Is(err, httpservice.TooManyRequestsError) {
		time.Sleep(backOffTimeout)

		// if we are enable to continue due to rate limiting then we should sleep for a bit
		// after the sleep we return the same link again so that it can be tried again
		t.links <- t.url
		return
	}

	if err != nil {
		log.Printf("Unable to crawl url = %s", t.url)
		return
	}

	for _, link := range links {
		t.links <- link
	}

	t.results <- TaskResult{url: t.url, links: links}
}
