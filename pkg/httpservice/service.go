package httpservice

import (
	"net/url"
)

type HttpService interface {
	GetHTML(path *url.URL) ([]byte, error)
}
