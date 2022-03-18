package linksservice

import "net/url"

type LinksParser interface {
	GetAllLinks(body []byte) []*url.URL
}
