package main

import (
	"time"

	"github.com/SlyMarbo/rss"
)

type Feed struct {
	url     string
	rawData *rss.Feed
}

func (feed *Feed) Fetch() {
	logger.Info("Stearted fetching field: %s", feed.url)

	rss.CacheParsedItemIDs(false)
	feed.rawData, _ = rss.Fetch(feed.url)

	logger.Info("Finished fetching field: %s", feed.url)
	logger.Info("There are %d items in %s", len(feed.rawData.Items), feed.url)
}

func (feed *Feed) ProcessNewItems(latest time.Time, action func(*FeedItem)) {
	for _, item := range feed.rawData.Items {
		feedItem := &FeedItem{*item}

		logger.Info("Checking item time: %v latest: %v", feedItem.Date.UTC(), latest.UTC())
		if feedItem.Date.UTC().After(latest.UTC()) {
			logger.Info("Enquing new item")
			action(feedItem)
		}
	}
}

type FeedItem struct {
	rss.Item
}

func (feed *FeedItem) NewArticle() *Article {
	return &Article{
		//ID:      feed.ID,
		Title:   feed.Title,
		Url:     feed.Link,
		Excerpt: feed.Summary,
		Date:    feed.Date,
	}
}
