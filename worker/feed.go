package worker

import (
	"time"

	"github.com/SlyMarbo/rss"
	. "github.com/rastasheep/utisak-worker/article"
	"github.com/rastasheep/utisak-worker/worker/parser"
)

type Feed struct {
	Url          string
	Category     string
	CategorySlug string `json:"category_slug"`
	Source       string
	SourceSlug   string `json:"source_slug"`
	Parser       string
	RawData      *rss.Feed
	Options      *parser.ParserOptions
}

func (feed *Feed) Fetch() error {
	var err error
	logger.Info("Started fetching feed: %s", feed.Url)

	rss.CacheParsedItemIDs(false)
	feed.RawData, err = rss.Fetch(feed.Url)
	if err != nil {
		return err
	}

	logger.Info("Finished fetching feed: %s", feed.Url)
	logger.Info("There are %d items in %s", len(feed.RawData.Items), feed.Url)
	return nil
}

func (feed *Feed) ProcessNewItems(latest time.Time, action func(*FeedItem)) {
	for _, item := range feed.RawData.Items {
		feedItem := &FeedItem{*item, *feed}

		logger.Info("Item time: %v latest article time for source: %v", feedItem.Date.UTC(), latest.UTC())
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
