package main

import (
	log "github.com/rastasheep/utisak-worker/log"
)

type FeedRegistry struct {
	Feeds  []*Feed
	Logger log.Logger
}

func NewFeedRegistry(sourcePath string) *FeedRegistry {
	var registry FeedRegistry

	registry.Logger = log.NewPrefixLogger("REGISTRY")
	LoadFile(config.FeedRegistryPath, &registry.Feeds)

	return &registry
}

func (registry *FeedRegistry) FetchFeeds(action func(*FeedItem)) {
	for _, feed := range registry.Feeds {
		feed.Fetch()

		article := feed.LatestArticle()

		registry.Logger.Info("Last article: ID: %b, Url: %s", article.ID, article.Url)
		feed.ProcessNewItems(article.Date, action)
	}
}
