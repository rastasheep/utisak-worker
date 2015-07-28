package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ParserAPI = "http://readability.com/api/content/v1/parser?url=%s&token=%s"

func ReadabilityParse(url string, target interface{}) error {
	apiUrl := fmt.Sprintf(ParserAPI, url, Conf.ReadabilityToken)
	r, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
