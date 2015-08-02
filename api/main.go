package main

import (
	"fmt"
	"net/http"

	log "github.com/rastasheep/utisak-worker/log"
)

var (
	logger log.Logger
)

func potsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("%s %s %s", r.RemoteAddr, r.Method, r.URL)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	if r.Method != "GET" {
		http.Error(w, fmt.Sprintf("Unsupported method '%s'\n", r.Method), 501)
		return
	}

	category := r.FormValue("category")
	if len(category) == 0 {
		category = "all"
	}

	fmt.Fprintf(w, "Hi there, I bring you  %s!", category)

}

func unknownHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("[404] %s %s %s", r.RemoteAddr, r.Method, r.URL)
	http.Redirect(w, r, "http://utisak.com", http.StatusFound)
}

func main() {
	log.LogTo("stdout", "DEBUG")
	logger = log.NewPrefixLogger("MAIN")

	http.HandleFunc("/posts", potsHandler)
	http.HandleFunc("/", unknownHandler)

	port := "8080" //os.Getenv("PORT")

	logger.Info("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
