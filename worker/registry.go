package main

import log "github.com/rastasheep/utisak-worker/log"

type FeedRegistry struct {
	feeds []*Feed
	log.Logger
}

func NewFeedRegistry(sourcePath string) *FeedRegistry {
	registry := &FeedRegistry{
		feeds:  make([]*Feed, 0),
		Logger: log.NewPrefixLogger("registry"),
	}
	registry.feeds = append(registry.feeds, &Feed{url: "http://www.politika.rs/rubrike/Sport/index.1.lt.xml"})
	return registry
}

func (registry *FeedRegistry) FetchFeeds(action func(*FeedItem)) {
	for _, feed := range registry.feeds {
		feed.Fetch()

		article := LatestArticle()

		logger.Info("Last article: ID: %b, Url: %s", article.ID, article.Url)
		feed.ProcessNewItems(article.Date, action)
	}
}
