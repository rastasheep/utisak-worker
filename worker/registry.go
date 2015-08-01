package main

import (
	"time"

	log "github.com/rastasheep/utisak-worker/log"
)

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
		//latestArticle =
		start, _ := time.Parse(time.RFC822, "24 Jul 15 18:00 UTC")
		feed.ProcessNewItems(start, action)
	}
}
