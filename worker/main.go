package main

import (
	"encoding/json"
	"fmt"

	log "github.com/rastasheep/utisak-worker/log"

	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	_ "github.com/lib/pq"
	"gopkg.in/robfig/cron.v2"
)

var (
	config       *Config
	logger       log.Logger
	db           *gorm.DB
	feedRegistry *FeedRegistry
)

func main() {
	config = LoadConfig()

	log.LogTo(config.LogTo, config.LogLevel)
	logger = log.NewPrefixLogger("MAIN")

	feedRegistry = NewFeedRegistry(config.FeedRegistryPath)

	db = newDb()
	defer db.Close()

	workers.Configure(config.RedisConfig())
	workers.Middleware.Append(&workers.MiddlewareRetry{})

	workers.Process("article_fetching", articleFetchingJob, 10)

	go workers.StatsServer(8080)

	go startCron()

	workers.Run()
}

func startCron() {
	c := cron.New()
	c.AddFunc("*/5 * * * * *", fetchFeeds)
	c.Start()
}

func fetchFeeds() {
	logger.Info("Starting to pull feeds")

	feedRegistry.FetchFeeds(enqueueArticleFetchingJob)

	log.Info("Finished pulling feeds")
}

func enqueueArticleFetchingJob(item *FeedItem) {
	workers.Enqueue("article_fetching", "Add", item)
}

func articleFetchingJob(message *workers.Msg) {
	var item FeedItem

	params := message.Args().ToJson()
	json.Unmarshal([]byte(params), &item)

	article := item.NewArticle()
	article.FetchDetails()

	logger.Info("Successfully created article: %+v\n", article)
}

func newDb() *gorm.DB {
	logger.Info("Connecting to postgres: %s", config.PostgresConfig())
	db, _ := gorm.Open("postgres", config.PostgresConfig())

	err := db.DB().Ping()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	db.AutoMigrate(&Article{})
	return &db
}
