package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/jrallison/go-workers"
	"gopkg.in/robfig/cron.v2"
)

var CronLogger = NewLogger("cron")
var JobLogger = NewLogger("job")
var Conf = LoadConfig()

func NewLogger(name string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%s: ", name), log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

var Urls = []string{
	"http://www.politika.rs/rubrike/Sport/index.1.lt.xml",
}

func main() {
	workers.Configure(map[string]string{
		"server":   Conf.Redis.Domain,
		"database": Conf.Redis.Database,
		"pool":     Conf.Redis.Pool,
		"process":  Conf.Redis.Process,
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
	CronLogger.Println("Starting to pull feeds")

	for _, url := range Urls {
		fetchFeed(url)
	}

	CronLogger.Println("Finished pulling feeds")
}

func fetchFeed(url string) {
	CronLogger.Printf("Stearted fetching field: %s\n", url)

	rss.CacheParsedItemIDs(false)
	feed, _ := rss.Fetch(url)

	CronLogger.Printf("Finished fetching field: %s\n", url)
	CronLogger.Printf("There are %d items in %s\n", len(feed.Items), url)

	fetchNewItems(feed.Items)
}

func fetchNewItems(items []*rss.Item) {
	start, _ := time.Parse(time.RFC822, "24 Jul 15 18:00 UTC")

	for _, item := range items {
		if item.Date.UTC().After(start) {
			CronLogger.Println("Enquing new item")
			workers.Enqueue("article_fetching", "Add", item)
		}
	}
}

func articleFetchingJob(message *workers.Msg) {
	params, _ := message.Args().Map()
	article := Article{}

	if err := article.ParseData(params); err != nil {
		panic(err)
	}

	JobLogger.Printf("Successfully created article: %+v\n", article)
}
