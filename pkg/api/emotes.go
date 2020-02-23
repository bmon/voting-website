package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type emote struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`
}

var emotesPath = "static/emotes/"
var emoteList []*emote

func (a *API) listEmotes(w http.ResponseWriter, r *http.Request) {
	if emoteList == nil {
		var err error
		emoteList, err = loadEmotes()
		if err != nil {

		}

	}

	writeJSON(w, emoteList)
	return
}

func loadEmotes() ([]*emote, error) {
	emotes := []*emote{}

	err := filepath.Walk(emotesPath, func(path string, info os.FileInfo, err error) error {
		filename := strings.TrimPrefix(path, emotesPath)
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		if name != "" {
			emotes = append(emotes, &emote{name, filename})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return emotes, nil
}
