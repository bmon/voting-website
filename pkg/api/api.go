package api

import (
	"fmt"
	"net/http"

	"github.com/bmon/voting-website/pkg/env"
)

type api struct {
	config *env.Config
}

func RegisterRoutes(config *env.Config) {
	a := &api{config}
	http.HandleFunc("/", a.Hello)
}

func (a *api) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
