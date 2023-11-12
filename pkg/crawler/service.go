package crawler

import "net/url"

type WebCrawler interface {
	Crawl(link *url.URL) ([]*url.URL, error)
}
