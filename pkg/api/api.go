package api

import (
	"context"
	"html/template"

	"github.com/bmon/voting-website/pkg/env"
	"github.com/coreos/go-oidc"
)

type API struct {
	Config        *env.Config
	indexTemplate *template.Template
	verifier      *oidc.IDTokenVerifier
}

func New() *API {
	config := env.LoadConfig()
	keySet := oidc.NewRemoteKeySet(context.Background(), "https://www.googleapis.com/oauth2/v3/certs")

	return &API{
		Config:        config,
		indexTemplate: template.Must(template.ParseFiles("index.html")),
		verifier: oidc.NewVerifier("https://accounts.google.com", keySet, &oidc.Config{
			ClientID: config.OauthClientID,
		}),
	}
}
