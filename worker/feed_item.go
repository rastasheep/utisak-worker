package worker

import (
	"encoding/json"
	"time"

	"github.com/SlyMarbo/rss"
	. "github.com/rastasheep/utisak-worker/article"
	"github.com/rastasheep/utisak-worker/worker/parser"
)

type FeedItem struct {
	rss.Item
	Feed
}

func (item *FeedItem) Fetch() {
	time.Sleep(3 * time.Second)

	article := item.newArticle()

	parser, err := parser.Get(item.Feed.Parser)
	if err != nil {
		logger.Error("FeedItem: unknown driver %q", item.Feed.Parser)
		return
	}

	if item.Feed.Options != nil {
		parser.SetOptions(*item.Feed.Options)
	}

	articleData, err := parser.Fetch(article.Url)
	if err != nil {
		logger.Error("FeedItem: error parsing article data %q", err)
		return
	}

	err = json.Unmarshal(articleData, article)
	if err != nil {
		logger.Error("FeedItem: failed to unmarshal article data %q", err)
		return
	}

	db.Create(&article)

	logger.Info("Successfully created article: %+v\n", article)
}

func (item *FeedItem) newArticle() *Article {
	return &Article{
		//ID:      feed.ID,
		Title:        item.Title,
		Url:          item.Link,
		Excerpt:      item.Summary,
		Date:         item.Date,
		Category:     item.Feed.Category,
		CategorySlug: item.Feed.CategorySlug,
		Source:       item.Feed.Source,
		SourceSlug:   item.Feed.SourceSlug,
	}
}
