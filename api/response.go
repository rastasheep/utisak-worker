package api

type Response struct {
	Articles   []SerializedArticle `json:"posts"`
	Categories []Category          `json:"categories"`
}
