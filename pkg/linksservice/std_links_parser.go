package linksservice

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/url"

	"golang.org/x/net/html"
)

const (
	anchorTag     = "a"
	hrefAttribute = "href"
)

type StdLinksParser struct {
}

func (lp *StdLinksParser) GetAllLinks(body []byte) []*url.URL {
	links := []*url.URL{}

	reader := html.NewTokenizer(bytes.NewReader(body))
	for {
		tokenType := reader.Next()

		switch tokenType {
		case html.ErrorToken:
			if errors.Is(reader.Err(), io.EOF) {
				return links
			}

			log.Printf("Error parsing HTML %s", reader.Err())
			continue

		case html.StartTagToken:
			if token := reader.Token(); token.Data == anchorTag {
				for _, attribute := range token.Attr {
					if attribute.Key != hrefAttribute {
						continue
					}

					link, err := url.Parse(attribute.Val)
					if err != nil {
						continue
					}

					links = append(links, link)
				}
			}
		}
	}
}

func NewLinksParser() LinksParser {
	return &StdLinksParser{}
}
