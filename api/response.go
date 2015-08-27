package api

import (
	. "github.com/rastasheep/utisak-worker/article"
)

type Response struct {
	Articles   []SerializedArticle `json:"posts"`
	Categories []Category          `json:"categories"`
}
