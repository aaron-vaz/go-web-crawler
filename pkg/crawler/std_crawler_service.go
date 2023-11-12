package crawler

import (
	"log"
	"net/url"

	"github.com/aaron-vaz/go-web-crawler/pkg/httpservice"
	"github.com/aaron-vaz/go-web-crawler/pkg/linksservice"
)

var (
	emptyLinks = make([]*url.URL, 0)
)

type StdWebCrawler struct {
	hs httpservice.HttpService
	ls linksservice.LinksParser
}

func (wc *StdWebCrawler) Crawl(link *url.URL) ([]*url.URL, error) {
	body, err := wc.hs.GetHTML(link)

	if err != nil {
		log.Printf("Error retrieving html from %s, error = %s", link, err)
		return emptyLinks, err
	}

	return wc.ls.GetAllLinks(body), nil
}

func NewWebCrawler(hs httpservice.HttpService, ls linksservice.LinksParser) WebCrawler {
	return &StdWebCrawler{hs: hs, ls: ls}
}
