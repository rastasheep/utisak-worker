package worker

import (
	"time"

	"github.com/SlyMarbo/rss"
	. "github.com/rastasheep/utisak-worker/article"
)

type FeedItem struct {
	rss.Item
	Feed
}

func (item *FeedItem) Fetch() {
	time.Sleep(3 * time.Second)

	article := item.newArticle()
	if err := FetchArticle(&item.Feed, article); err != nil {
		logger.Info("Unable to fetch article: %s", err)
	}

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
