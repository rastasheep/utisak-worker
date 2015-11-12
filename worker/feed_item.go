package worker

import (
	"time"

	"github.com/SlyMarbo/rss"
	. "github.com/rastasheep/utisak-worker/article"
)

type FeedItem struct {
	rss.Item
	Category     string
	CategorySlug string
	Source       string
	SourceSlug   string
}

func (item *FeedItem) Fetch() {
	time.Sleep(3 * time.Second)

	article := item.newArticle()
	ReadabilityParse(article.Url, &article)

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
