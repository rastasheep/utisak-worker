package parser

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const ReadabilityAPI = "http://readability.com/api/content/v1/parser"

type readability struct{}

func (parser *readability) Fetch(sourceUrl string) ([]byte, error) {
	fullUrl, _ := url.Parse(ReadabilityAPI)
	parameters := url.Values{}
	parameters.Add("url", sourceUrl)
	parameters.Add("token", os.Getenv("READABILITY_PASSWORD"))
	fullUrl.RawQuery = parameters.Encode()

	r, err := http.Get(fullUrl.String())
	defer r.Body.Close()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r.Body)
}

func init() {
	RegisterParser("readability", &readability{})
}
