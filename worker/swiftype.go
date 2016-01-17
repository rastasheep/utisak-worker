package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"

	. "github.com/rastasheep/utisak-worker/article"
)

const IndexAPI = "https://api.swiftype.com/api/v1/engines/%s/document_types/%s/documents/bulk_create_or_update_verbose.json"

type stBody struct {
	AuthToken string        `json:"auth_token"`
	Documents []*stDocument `json:"documents"`
}

type stDocument struct {
	ExternalId string    `json:"external_id"`
	Fields     []stField `json:"fields"`
}

type stField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func StIndexArticles(articles []SerializedArticle) error {
	articlesSize := len(articles)
	stDocuments := make([]*stDocument, articlesSize)

	if articlesSize == 0 {
		return nil
	}

	for i := 0; i < articlesSize; i++ {
		article := articles[i]
		stDocuments[i] = &stDocument{
			ExternalId: fmt.Sprintf("%d", article.ID),
			Fields: []stField{
				stField{Name: "title", Value: article.Title, Type: "string"},
				stField{Name: "domain", Value: article.Domain, Type: "enum"},
				stField{Name: "url", Value: article.ShareUrl, Type: "enum"},
				stField{Name: "excerpt", Value: html.UnescapeString(article.Excerpt), Type: "text"},
				stField{Name: "word_count", Value: fmt.Sprintf("%d", article.WordCount), Type: "integer"},
				stField{Name: "published_at", Value: fmt.Sprintf("%s", article.Date), Type: "date"},
				stField{Name: "hero_image_url", Value: article.LeadImage, Type: "enum"},
				stField{Name: "category", Value: article.Category, Type: "string"},
				stField{Name: "category_slug", Value: article.CategorySlug, Type: "enum"},
				stField{Name: "author", Value: article.Source, Type: "enum"},
				stField{Name: "total_views", Value: fmt.Sprintf("%d", article.TotalViews), Type: "integer"},
			},
		}
	}

	return indexStDocuments(stDocuments)
}

func indexStDocuments(stDocuments []*stDocument) error {
	conf := config.SwiftypeConfig()

	indexUrl := fmt.Sprintf(IndexAPI, conf.Engine, conf.DocumentType)

	body := &stBody{
		AuthToken: conf.AuthToken,
		Documents: stDocuments,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("Unable to marshal swiftype body: %s", jsonBody))
	}

	buff := bytes.NewBuffer(jsonBody)
	resp, err := http.Post(indexUrl, "application/json", buff)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	logger.Info("[ST] Response: %+v - %+v", string(resp.Status), string(respBody))
	logger.Info("[ST] Successfully indexed: %+v\n", body)

	return nil
}
