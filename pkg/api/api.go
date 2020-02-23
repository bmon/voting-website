package api

import (
	"context"
	"html/template"
	"log"

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
