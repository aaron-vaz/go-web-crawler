package httpservice

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	NoBodyError          = errors.New("Failed to retrive HTTP body")
	TooManyRequestsError = errors.New("Too many requests to host")
)

var (
	emptyBytes = make([]byte, 0)
)

type StdHttpService struct {
	client *http.Client
}

func (gs *StdHttpService) GetHTML(path *url.URL) ([]byte, error) {
	res, err := gs.client.Get(path.String())
	if err != nil {
		return emptyBytes, err
	}

	if res.StatusCode == http.StatusTooManyRequests {
		return emptyBytes, TooManyRequestsError
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Error making http request = %s, code = %d", path, res.StatusCode)
		return emptyBytes, NoBodyError
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return emptyBytes, NoBodyError
	}

	return body, nil
}

func NewHttpService(client *http.Client) HttpService {
	return &StdHttpService{client: client}
}
