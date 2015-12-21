package worker

import (
	"fmt"

	. "github.com/rastasheep/utisak-worker/article"
	log "github.com/rastasheep/utisak-worker/log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/robfig/cron.v2"
)

var (
	// commit sha for the current build, set by the compile process.
	version  string
	revision string

	config *Config
	logger log.Logger
	db     *gorm.DB
)

const indexBatchSize = 20
const reFetchBatchSize = 20

func Main() {
	config = LoadConfig()
	BaseUrl = config.BaseUrl
	ArticlePrefix = config.ArticlePrefix

	log.LogTo(config.LogTo, config.LogLevel)
	logger = log.NewPrefixLogger("MAIN")

	db = newDb()
	defer db.Close()

	startCron()
	wait()
}

func startCron() {
	c := cron.New()
	c.AddFunc("0 */5 * * * *", fetchFeeds)
	c.AddFunc("0 */6 * * * *", reFetchFeeds)

	if config.Swiftype.Enabled {
		c.AddFunc("0 */6 * * * *", indexArticles)
	}

	c.Start()
}

func fetchFeeds() {
	logger.Info("Starting to pull feeds")

	worker := NewBackgroundWorker(5)
	fetchingJob := func(item *FeedItem) { worker.Queue <- item.Fetch }

	feedRegistry := NewFeedRegistry(config.FeedRegistryPath)
	feedRegistry.FetchFeeds(fetchingJob)

	worker.Process()
	log.Info("Finished pulling feeds")
}

func reFetchFeeds() {
	logger.Info("[RF] Starting to re-fetching feeds")

	var articles []Article
	feedRegistry := NewFeedRegistry(config.FeedRegistryPath)

	db.Where("refetch").Order("date desc").Limit(reFetchBatchSize).Find(&articles)

	for _, article := range articles {
		feed, err := feedRegistry.FindFeed(article.SourceSlug, article.CategorySlug)

		if err != nil {
			logger.Info("[RF] %s", err)
			continue
		}

		if err := FetchArticle(feed, &article); err != nil {
			logger.Info("[RF] Unable to fetch article: %s", err)
		}

		db.Model(&article).UpdateColumn("refetch", "true")
	}

	log.Info("[RF] Finished re-fetching feeds")
}

func indexArticles() {
	logger.Info("[ST] Starting to index articles")

	var articles []SerializedArticle
	db.Where("NOT indexed").Order("date desc").Limit(indexBatchSize).Find(&articles)

	if err := StIndexArticles(articles); err != nil {
		logger.Info("[ST] Unable to post swiftype documents: %s", err)
	}

	for _, article := range articles {
		db.Model(&article).Where("id = ?", article.ID).UpdateColumn("indexed", "true")
	}
	logger.Info("[ST] Finished indexing articles")
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

// blocks forever
func wait() {
	var ch chan bool
	<-ch
}
