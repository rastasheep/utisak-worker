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

		if feedItem.Date.UTC().After(latest) {
			logger.Info("Enquing new item")
			action(feedItem)
			logger.Info("article_fetching", "Add", feedItem)
		}
	}
}

type FeedItem struct {
	rss.Item
}

func (feed *FeedItem) NewArticle() *Article {
	return &Article{
		ID:      feed.ID,
		Title:   feed.Title,
		Url:     feed.Link,
		Excerpt: feed.Summary,
		Date:    feed.Date,
	}
}
