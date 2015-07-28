package main

import (
	"time"

	log "github.com/rastasheep/utisak-worker/log"

	"github.com/SlyMarbo/rss"
	"github.com/jrallison/go-workers"
	"gopkg.in/robfig/cron.v2"
)

var (
	config       *Config
	logger       log.Logger
	feedRegistry *FeedRegistry
	Urls         = []string{
		"http://www.politika.rs/rubrike/Sport/index.1.lt.xml",
	}
)

func main() {
	log.LogTo("stdout", "DEBUG")
	logger = log.NewPrefixLogger("MAIN")

	config = LoadConfig()
	feedRegistry = NewFeedRegistry(config.FeedRegistryPath)

	workers.Configure(map[string]string{
		"server":   config.Redis.Domain,
		"database": config.Redis.Database,
		"pool":     config.Redis.Pool,
		"process":  config.Redis.Process,
	})

	workers.Middleware.Append(&workers.MiddlewareRetry{})

	workers.Process("article_fetching", articleFetchingJob, 10)

	go workers.StatsServer(8080)

	go startCron()

	workers.Run()
}

func startCron() {
	c := cron.New()
	c.AddFunc("*/5 * * * * *", pullFeeds)
	c.Start()
}

func pullFeeds() {
	logger.Info("Starting to pull feeds")

	for _, url := range Urls {
		fetchFeed(url)
	}

	log.Info("Finished pulling feeds")
}

func fetchFeed(url string) {
	logger.Info("Stearted fetching field: %s", url)

	rss.CacheParsedItemIDs(false)
	feed, _ := rss.Fetch(url)

	logger.Info("Finished fetching field: %s", url)
	logger.Info("There are %d items in %s", len(feed.Items), url)

	fetchNewItems(feed.Items)
}

func fetchNewItems(items []*rss.Item) {
	start, _ := time.Parse(time.RFC822, "24 Jul 15 18:00 UTC")

	for _, item := range items {
		if item.Date.UTC().After(start) {
			logger.Info("Enquing new item")
			logger.Info("article_fetching", "Add", item)
		}
	}
}

func articleFetchingJob(message *workers.Msg) {
	params, _ := message.Args().Map()
	article := Article{}

	if err := article.ParseData(params); err != nil {
		panic(err)
	}

	logger.Info("Successfully created article: %+v\n", article)
}
