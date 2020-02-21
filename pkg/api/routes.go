package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *API) Router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/hello", a.hello)
	r.HandleFunc("/whoami", a.mustVerify(a.whoami))
	r.HandleFunc("/", a.index)

	return r
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{error: "%s"}`, msg)
	return
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
