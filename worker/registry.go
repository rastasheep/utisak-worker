package main

import (
	log "github.com/rastasheep/utisak-worker/log"
)

type FeedRegistry struct {
	feeds []*Feed
	log.Logger
}

type Feed struct {
}

func NewFeedRegistry(sourcePath string) *FeedRegistry {
	registry := &FeedRegistry{
		feeds:  make([]*Feed, 0),
		Logger: log.NewPrefixLogger("registry"),
	}
	return registry

}
