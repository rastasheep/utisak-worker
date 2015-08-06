package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	. "github.com/rastasheep/utisak-worker/article"
	log "github.com/rastasheep/utisak-worker/log"
)

var (
	config *Config
	logger log.Logger
	db     *gorm.DB
)

func potsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("%s %s %s", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Unsupported method '%s'\n", r.Method), 501)
		return
	}

	var articles []SerializedArticle

	category := strings.Split(r.FormValue("category"), ",")
	category = deleteEmpty(category)

	if len(category) == 0 || category[0] == "all" {
		db.Find(&articles)
	} else {
		db.Where("category_slug IN (?)", category).Find(&articles)
	}

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

	http.HandleFunc("/posts", serve(potsHandler))
	http.HandleFunc("/", unknownHandler)

	port := "8080" //os.Getenv("PORT")

	logger.Info("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
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
