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
	Category     string
	CategorySlug string
	Source       string
	SourceSlug   string
	Parser       string
}

func (item *FeedItem) Fetch() {
	time.Sleep(3 * time.Second)

	article := item.newArticle()

	parser, err := parser.Get(item.Parser)
	if err != nil {
		logger.Error("FeedItem: unknown driver %q", item.Parser)
		return
	}

	articleData, err := parser.Parse(article.Url)
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

func (feed *FeedItem) newArticle() *Article {
	return &Article{
		//ID:      feed.ID,
		Title:        feed.Title,
		Url:          feed.Link,
		Excerpt:      feed.Summary,
		Date:         feed.Date,
		Category:     feed.Category,
		CategorySlug: feed.CategorySlug,
		Source:       feed.Source,
		SourceSlug:   feed.SourceSlug,
	}
}
