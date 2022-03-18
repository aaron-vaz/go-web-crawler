package worker

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/aaron-vaz/go-web-crawler/pkg/crawler"
)

const (
	idleTimeout = 5 * time.Second
)

var (
	visitedLinks sync.Map
)

func resetTimer(timer *time.Timer) {
	timer.Stop()
	timer.Reset(idleTimeout)
}

type TaskDispactcher struct {
	startURL *url.URL
	wc       crawler.WebCrawler
}

func (td *TaskDispactcher) Start() []TaskResult {
	results := []TaskResult{}

	links := make(chan *url.URL, 100)
	taskResults := make(chan TaskResult)

	td.runTaskAsync(td.startURL, links, taskResults)

	idleTimer := time.NewTimer(idleTimeout)
	for {
		select {

		case link := <-links:
			filteredLink := td.filterLink(link)
			if filteredLink == nil {
				continue
			}

			td.runTaskAsync(filteredLink, links, taskResults)

			resetTimer(idleTimer)

		case result := <-taskResults:
			visitedLinks.Store(result.url, result.url)
			results = append(results, result)

		case <-idleTimer.C:
			log.Println("Dispatcher Idle timeout reached, stoping dispatcher")
			return results
		}

	}
}

func (td *TaskDispactcher) filterLink(link *url.URL) *url.URL {
	var resolvedLink = link
	if link.Host == "" {
		resolvedLink = td.startURL.ResolveReference(link)
	}

	_, visited := visitedLinks.Load(resolvedLink)
	if visited || resolvedLink.Host != td.startURL.Host {
		return nil
	}

	return resolvedLink
}

func (td *TaskDispactcher) runTaskAsync(url *url.URL, links chan *url.URL, results chan TaskResult) {
	startURLTask := Task{
		url:     url,
		links:   links,
		results: results,

		wc: td.wc,
	}

	go startURLTask.Run()
}

func NewTaskDispactcher(startURL *url.URL, wc crawler.WebCrawler) *TaskDispactcher {
	return &TaskDispactcher{startURL: startURL, wc: wc}
}
