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
	startURL   *url.URL
	maxWorkers int

	wc crawler.WebCrawler
}

func (td *TaskDispactcher) Start() []TaskResult {
	results := []TaskResult{}

	links := make(chan *url.URL)
	taskResults := make(chan TaskResult)

	td.startTasks(links, taskResults)

	links <- td.startURL

	idleTimer := time.NewTimer(idleTimeout)
	for {
		select {

		case taskResult := <-taskResults:
			resetTimer(idleTimer)

			log.Printf("Got result = %+v", taskResult)

			if taskResult.retry {
				links <- taskResult.url
				break
			}

			visitedLinks.Store(taskResult.url.String(), taskResult.url)
			results = append(results, taskResult)

			for _, link := range taskResult.links {
				filteredLink := td.filterLink(link)
				if filteredLink == nil {
					continue
				}

				links <- filteredLink
			}

		case <-idleTimer.C:
			log.Println("Dispatcher Idle timeout reached, stoping dispatcher")

			close(links)
			close(taskResults)

			return results
		}

	}
}

func (td *TaskDispactcher) filterLink(link *url.URL) *url.URL {
	var resolvedLink = link
	if link.Host == "" {
		resolvedLink = td.startURL.ResolveReference(link)
	}

	_, visited := visitedLinks.Load(resolvedLink.String())
	if visited || resolvedLink.Host != td.startURL.Host {
		return nil
	}

	return resolvedLink
}

func (td *TaskDispactcher) startTasks(links chan *url.URL, results chan TaskResult) {
	for count := 0; count <= td.maxWorkers; count++ {
		task := Task{links: links, results: results, wc: td.wc}
		go task.Run()
	}
}

func NewTaskDispactcher(startURL *url.URL, maxWorkers int, wc crawler.WebCrawler) *TaskDispactcher {
	return &TaskDispactcher{startURL: startURL, maxWorkers: maxWorkers, wc: wc}
}
