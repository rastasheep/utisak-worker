package worker

import (
	"fmt"

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
		err := feed.Fetch()
		if err != nil {
			registry.Logger.Error("Failed to fetch feed: %s reason: %s", feed.Url, err.Error())
			continue
		}

		article := feed.LatestArticle()

		registry.Logger.Info("Last article: ID: %b, Url: %s", article.ID, article.Url)
		feed.ProcessNewItems(article.Date, action)
	}
}

func (registry *FeedRegistry) FindFeed(sourceSlug, categorySlug string) (*Feed, error) {
	for _, feed := range registry.Feeds {
		if feed.SourceSlug == sourceSlug && feed.CategorySlug == categorySlug {
			return feed, nil
		}
	}
	return nil, fmt.Errorf("No feeds found for: '%q' '%q'", sourceSlug, categorySlug)
}
