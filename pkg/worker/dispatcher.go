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
	visitedLinks        sync.Map
	currentTargetWorker int
)

func resetTimer(timer *time.Timer) {
	timer.Stop()
	timer.Reset(idleTimeout)
}

type TaskDispactcher struct {
	startURL       *url.URL
	workerChannels []chan *url.URL
	maxWorkers     int

	wc crawler.WebCrawler
}

func (td *TaskDispactcher) Start() []TaskResult {
	results := make([]TaskResult, 0, td.maxWorkers)
	taskResults := make(chan TaskResult, td.maxWorkers)

	td.startTasks(taskResults)
	td.publishToTask(td.startURL)

	idleTimer := time.NewTimer(idleTimeout)
	for {
		select {
		case taskResult := <-taskResults:
			resetTimer(idleTimer)

			if taskResult.retry {
				td.publishToTask(taskResult.url)
				break
			}

			visitedLinks.Store(taskResult.url.String(), taskResult.url)
			results = append(results, taskResult)

			for _, link := range taskResult.links {
				filteredLink := td.filterLink(link)
				if filteredLink == nil {
					continue
				}

				td.publishToTask(filteredLink)
			}

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

	_, visited := visitedLinks.Load(resolvedLink.String())
	if visited || resolvedLink.Host != td.startURL.Host {
		return nil
	}

	return resolvedLink
}

func (td *TaskDispactcher) startTasks(results chan TaskResult) {
	for count := 0; count < td.maxWorkers; count++ {
		links := make(chan *url.URL)
		td.workerChannels[count] = links

		task := Task{links: links, results: results, wc: td.wc}
		task.Run()
	}
}

func (td *TaskDispactcher) publishToTask(link *url.URL) {
	workerChannel := td.workerChannels[currentTargetWorker]

	go func(link *url.URL, workerChannel chan *url.URL) {
		workerChannel <- link
	}(link, workerChannel)

	if currentTargetWorker++; currentTargetWorker > len(workerChannel)-1 {
		currentTargetWorker = 0
	}
}

func NewTaskDispactcher(startURL *url.URL, maxWorkers int, wc crawler.WebCrawler) *TaskDispactcher {
	return &TaskDispactcher{
		startURL:       startURL,
		workerChannels: make([]chan *url.URL, maxWorkers),
		maxWorkers:     maxWorkers,
		wc:             wc,
	}
}
