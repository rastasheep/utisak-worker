package api

type Response struct {
	Categories []Category `json:"categories"`
	Articles   []Article  `json:"posts"`
}
