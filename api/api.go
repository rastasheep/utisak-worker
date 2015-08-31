package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	. "github.com/rastasheep/utisak-worker/article"
	log "github.com/rastasheep/utisak-worker/log"
)

const (
	perPage = 20
	allCat  = "all"
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
	articleID := mux.Vars(r)["article_id"]
	logger.Info("[article] %s ", articleID)

	var article SerializedArticle

	db.Select("url").Limit(1).Find(&article, articleID)

	go db.Model(&article).Where("id = ?", articleID).UpdateColumn("total_views", gorm.Expr("total_views + ?", 1))

	http.Redirect(w, r, article.Url, 301)
}
func articlesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	category := strings.Split(r.FormValue("category"), ",")
	category = deleteEmpty(category)

	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		page = 0
	}

	sort := r.FormValue("sort")

	var articles []SerializedArticle
	searchRelation := db

	if len(category) != 0 && !Contains(category, allCat) {
		searchRelation = searchRelation.Where("category_slug IN (?)", category)
	}

	if page > 1 {
		searchRelation = searchRelation.Offset((page - 1) * perPage)
	}

	if sort == "newest" {
		searchRelation = searchRelation.Order("date desc")
	} else {
		searchRelation = searchRelation.Order("(total_views / POW(((EXTRACT(EPOCH FROM (now()-date)) / 3600)::integer + 2), 1.5)) desc")
	}

	searchRelation.Limit(perPage).Find(&articles)

	categories := GetCategories()

	if resp, err := json.Marshal(&Response{Categories: categories, Articles: articles}); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func unknownHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("[404] %s %s %s", r.RemoteAddr, r.Method, r.URL)
	http.Redirect(w, r, "http://utisak.com", http.StatusFound)
}

func Main() {
	config = LoadConfig()

	log.LogTo(config.LogTo, config.LogLevel)
	logger = log.NewPrefixLogger("MAIN")

	db = newDb()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/posts/{article_id}", serve(articleHandler)).Methods("GET")
	r.HandleFunc("/posts", serve(articlesHandler)).Methods("GET")
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

func deleteEmpty(s []string) []string {
	var r []string
	if len(s) == 0 {
		return s
	}
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func Contains(slice interface{}, val interface{}) bool {
	sv := reflect.ValueOf(slice)

	for i := 0; i < sv.Len(); i++ {
		if sv.Index(i).Interface() == val {
			return true
		}
	}
	return false
}
