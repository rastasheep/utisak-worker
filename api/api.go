package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	. "github.com/rastasheep/utisak-worker/article"
	log "github.com/rastasheep/utisak-worker/log"
)

var (
	// commit sha for the current build, set by the compile process.
	version  string
	revision string

	config *Config
	logger log.Logger
	db     *gorm.DB
)

func articleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceSlug := vars["source_slug"]
	articleSlug := vars["article_slug"]
	articleID := vars["article_id"]
	logger.Info("[article] Find for: source_slug=%s, article_slug:%s, article_id:%s", sourceSlug, articleSlug, articleID)

	var article SerializedArticle

	if db.Select("url").Where(map[string]interface{}{"id": articleID, "slug": articleSlug, "source_slug": sourceSlug}).Limit(1).Find(&article).RecordNotFound() {
		logger.Info("[article] Not found: source_slug=%s, article_slug:%s, article_id:%s", sourceSlug, articleSlug, articleID)
		unknownHandler(w, r)
		return
	}

	go db.Model(&article).Where(map[string]interface{}{"id": articleID, "slug": articleSlug, "source_slug": sourceSlug}).UpdateColumn("total_views", gorm.Expr("total_views + ?", 1))

	http.Redirect(w, r, article.Url, 301)
}

func reFetchArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sourceSlug := vars["source_slug"]
	articleSlug := vars["article_slug"]
	articleID := vars["article_id"]
	logger.Info("[article] Find for: source_slug=%s, article_slug:%s, article_id:%s", sourceSlug, articleSlug, articleID)

	var article SerializedArticle

	if db.Select("url").Where(map[string]interface{}{"id": articleID, "slug": articleSlug, "source_slug": sourceSlug}).Limit(1).Find(&article).RecordNotFound() {
		logger.Info("[article] Not found: source_slug=%s, article_slug:%s, article_id:%s", sourceSlug, articleSlug, articleID)
		unknownHandler(w, r)
		return
	}

	go db.Model(&article).Where(map[string]interface{}{"id": articleID, "slug": articleSlug, "source_slug": sourceSlug}).UpdateColumn("refetch", "true")

	w.WriteHeader(200)
}

func articlesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	var articles []SerializedArticle

	params := &QueryParams{
		CategoryStr:  r.FormValue("category"),
		PageStr:      r.FormValue("page"),
		Sort:         r.FormValue("sort"),
		StartDateStr: r.FormValue("start-date"),
		EndDateStr:   r.FormValue("end-date"),
	}
	params.Parse()

	params.PrepareQuery(db).Find(&articles)
	categories := GetCategories()

	if resp, err := json.Marshal(&Response{Categories: categories, Articles: articles}); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func unknownHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("[404] %s %s %s", r.RemoteAddr, r.Method, r.URL)
	http.Redirect(w, r, BaseUrl, http.StatusFound)
}

func Main() {
	config = LoadConfig()
	BaseUrl = config.BaseUrl
	ArticlePrefix = config.ArticlePrefix

	log.LogTo(config.LogTo, config.LogLevel)
	logger = log.NewPrefixLogger("MAIN")

	db = newDb()
	defer db.Close()

	artclesPath := fmt.Sprintf("/%s", ArticlePrefix)
	articlePath := fmt.Sprintf("/%s/{source_slug}/{article_slug}/{article_id:[0-9]+}", ArticlePrefix)

	r := mux.NewRouter()
	r.HandleFunc(articlePath, serve(articleHandler)).Methods("GET")
	r.HandleFunc(articlePath, serve(reFetchArticleHandler)).Methods("PUT")
	r.HandleFunc(artclesPath, serve(articlesHandler)).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(unknownHandler)

	port := "8080" //os.Getenv("PORT")

	logger.Info("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, r)
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

func serve(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fn(w, r) // call the handler

		elapsed := time.Now().Sub(start)

		logger.Info("%v %s %s", elapsed, r.Method, r.RequestURI)
	}
}
