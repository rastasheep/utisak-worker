package worker

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/rastasheep/utisak-worker/article"
	log "github.com/rastasheep/utisak-worker/log"

	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	_ "github.com/lib/pq"
	"gopkg.in/robfig/cron.v2"
)

var (
	// commit sha for the current build, set by the compile process.
	version  string
	revision string

	config       *Config
	logger       log.Logger
	db           *gorm.DB
	feedRegistry *FeedRegistry
)

func Main() {
	config = LoadConfig()
	BaseUrl = config.BaseUrl
	ArticlePrefix = config.ArticlePrefix

	log.LogTo(config.LogTo, config.LogLevel)
	logger = log.NewPrefixLogger("MAIN")

	feedRegistry = NewFeedRegistry(config.FeedRegistryPath)
	logger.Info("Registry: %+v", feedRegistry)

	db = newDb()
	defer db.Close()

	workers.Configure(config.RedisConfig())
	workers.Middleware.Append(&workers.MiddlewareRetry{})

	workers.Process("article_fetching", articleFetchingJob, 5)
	workers.Process("search_indexing", searchIndexingJob, 5)

	go workers.StatsServer(8080)

	go startCron()

	workers.Run()
}

func startCron() {
	c := cron.New()
	c.AddFunc("0 */5 * * * *", fetchFeeds)
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

	time.Sleep(5 * time.Second)

	params := message.Args().ToJson()
	json.Unmarshal([]byte(params), &item)

	article := item.NewArticle()
	ReadabilityParse(article.Url, &article)

	db.Create(&article)

	if config.Swiftype.Enabled {
		logger.Info("[st] Enquing search indexing job")
		enqueueSearchIndexingJob(article.ID)
	}

	logger.Info("Successfully created article: %+v\n", article)
}

func enqueueSearchIndexingJob(id uint) {
	workers.Enqueue("search_indexing", "Add", id)
}

func searchIndexingJob(message *workers.Msg) {
	var article SerializedArticle
	id, _ := message.Args().Uint64()
	logger.Info("[st] Indexing job received id: %d", id)

	if db.First(&article, id).RecordNotFound() {
		panic(fmt.Sprintf("Unable to find article for indexing: %s", id))
	}

	if err := SwiftypeIndex(&article); err != nil {
		panic(fmt.Sprintf("Unable to post swiftype document: %s", err))
	}
}

func newDb() *gorm.DB {
	logger.Info("Connecting to postgres: %s", config.PostgresConfig())
	db, _ := gorm.Open("postgres", config.PostgresConfig())

	err := db.DB().Ping()
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	db.LogMode(true)
	db.AutoMigrate(&Article{})
	return &db
}
