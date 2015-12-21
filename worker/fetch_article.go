package worker

import (
	"encoding/json"
	"fmt"

	. "github.com/rastasheep/utisak-worker/article"
	"github.com/rastasheep/utisak-worker/worker/parser"
)

func FetchArticle(feed *Feed, article *Article) error {
	parser, err := parser.Get(feed.Parser)
	if err != nil {
		return fmt.Errorf("FeedItem: unknown driver %q", feed.Parser)
	}

	if feed.Options != nil {
		parser.SetOptions(*feed.Options)
	}

	articleData, err := parser.Fetch(article.Url)
	if err != nil {
		return fmt.Errorf("FeedItem: error parsing article data %q", err)
	}

	err = json.Unmarshal(articleData, article)
	if err != nil {
		return fmt.Errorf("FeedItem: failed to unmarshal article data %q", err)
	}

	if db.NewRecord(article) {
		db.Create(&article)
		return nil
	}
	db.Save(&article)
	return nil
}
