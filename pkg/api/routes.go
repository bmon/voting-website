package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (a *API) Router() http.Handler {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/hello", a.hello)
	r.HandleFunc("/whoami", a.mustVerify(a.whoami))
	r.HandleFunc("/emotes", a.mustVerify(a.listEmotes))
	r.HandleFunc("/votes", a.mustVerify(a.listVotes))
	r.HandleFunc("/vote", a.mustVerify(a.voteEmote))
	r.HandleFunc("/", a.index)

	// log to stdout
	return handlers.LoggingHandler(os.Stdout, r)
}

func (a *API) hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (a *API) index(w http.ResponseWriter, r *http.Request) {
	err := a.indexTemplate.Execute(w, a)
	if err != nil {
		fmt.Println(err)
	}
}

func (a *API) whoami(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("user").(*User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, u)
}
