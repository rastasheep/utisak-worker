package api

import . "github.com/rastasheep/utisak-worker/article"

type Response struct {
	Articles   []Article  `json:"posts"`
	Categories []Category `json:"categories"`
}
