package api

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/bmon/voting-website/pkg/env"
	"github.com/coreos/go-oidc"
)

type API struct {
	Config        *env.Config
	indexTemplate *template.Template
	verifier      *oidc.IDTokenVerifier
	store         *firestore.Client
}

func New() *API {
	config := env.LoadConfig()
	keySet := oidc.NewRemoteKeySet(context.Background(), "https://www.googleapis.com/oauth2/v3/certs")
	fbclient, err := firestore.NewClient(context.Background(), config.ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	return &API{
		Config:        config,
		indexTemplate: template.Must(template.ParseFiles("index.html")),
		verifier: oidc.NewVerifier("https://accounts.google.com", keySet, &oidc.Config{
			ClientID: config.OauthClientID,
		}),
		store: fbclient,
	}
}

func (a *API) Shutdown() {
	a.store.Close()
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
