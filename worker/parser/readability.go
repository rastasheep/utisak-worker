package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const ParserAPI = "http://readability.com/api/content/v1/parser?url=%s&token=%s"

type readability struct{}

func (parser *readability) Parse(url string, target interface{}) error {
	apiUrl := fmt.Sprintf(ParserAPI, url, os.Getenv("READABILITY_PASSWORD"))
	r, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func init() {
	RegisterParser("readability", &readability{})
}
