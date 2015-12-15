package parser

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const UtisakParserAPI = "http://parser.utisak.com/v1"

type copyist struct {
	Options ParserOptions
}

type CopyistImage struct {
	Type   string
	Url    string
	Height int
	Width  int
}

type CopyistArticle struct {
	LeadImage string `json:"lead_image_url"`
	Excerpt   string
	WordCount int `json:"word_count"`
	Domain    string

	// private
	Text  string       `json:"text"`
	Image CopyistImage `json:"image"`
}

func (article *CopyistArticle) populate() {
	article.LeadImage = article.Image.Url
	article.Text = strings.Trim(article.Text, "\n ")
	article.Text = strings.Replace(article.Text, "\n", " ", -1)
	article.Text = strings.Replace(article.Text, "  ", " ", -1)

	article.Excerpt = article.Text
	if len(article.Text) > 250 {
		article.Excerpt = article.Text[:250] + "..."
	}

	words := strings.Fields(article.Text)
	article.WordCount = len(words)

	return
}

func (parser *copyist) Fetch(sourceUrl string) ([]byte, error) {
	fullUrl, _ := url.Parse(UtisakParserAPI)
	parameters := url.Values{}
	parameters.Add("url", sourceUrl)
	if parser.Options.Language != "" {
		parameters.Add("lang", parser.Options.Language)
	}
	fullUrl.RawQuery = parameters.Encode()

	r, err := http.Get(fullUrl.String())
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	article := CopyistArticle{}

	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		return nil, err
	}
	article.populate()

	return json.Marshal(article)
}

func (parser *copyist) SetOptions(options ParserOptions) {
	parser.Options = options
}

func init() {
	RegisterParser("utisak-parser", &copyist{})
}
