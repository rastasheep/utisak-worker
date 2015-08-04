package worker

import (
	"time"

	"github.com/SlyMarbo/rss"
	. "github.com/rastasheep/utisak-worker/article"
)

type Feed struct {
	Url          string
	Category     string
	CategorySlug string `json:"category_slug"`
	Source       string
	RawData      *rss.Feed
}

func (feed *Feed) Fetch() {
	logger.Info("Stearted fetching field: %s", feed.Url)

	rss.CacheParsedItemIDs(false)
	feed.RawData, _ = rss.Fetch(feed.Url)

	logger.Info("Finished fetching field: %s", feed.Url)
	logger.Info("There are %d items in %s", len(feed.RawData.Items), feed.Url)
}

func (feed *Feed) ProcessNewItems(latest time.Time, action func(*FeedItem)) {
	for _, item := range feed.RawData.Items {
		feedItem := &FeedItem{*item, feed.Category, feed.CategorySlug, feed.Source}

		logger.Info("Checking item time: %v latest: %v", feedItem.Date.UTC(), latest.UTC())
		if feedItem.Date.UTC().After(latest.UTC()) {
			logger.Info("Enquing new item")
			action(feedItem)
		}
	}
}

func (feed *Feed) LatestArticle() *Article {
	var article Article

	db.Where("category_slug = ? and source = ?", feed.CategorySlug, feed.Source).
		Order("date desc").
		Limit(1).
		Find(&article)

	return &article
}

type FeedItem struct {
	rss.Item
	Category     string
	CategorySlug string
	Source       string
}

func (feed *FeedItem) NewArticle() *Article {
	return &Article{
		//ID:      feed.ID,
		Title:        feed.Title,
		Url:          feed.Link,
		Excerpt:      feed.Summary,
		Date:         feed.Date,
		Category:     feed.Category,
		CategorySlug: feed.CategorySlug,
		Source:       feed.Source,
	}
}
