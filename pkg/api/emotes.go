package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Emote struct {
	Name     string `json:"name" firestore:"name"`
	Filename string `json:"filename" firestore:"-"`
	Count    int64  `firestore:"count"`
}

var emotesPath = "static/emotes/"
var emoteList []*Emote

func (a *API) listEmotes(w http.ResponseWriter, r *http.Request) {
	if emoteList == nil {
		var err error
		emoteList, err = loadEmotes()
		if err != nil {
			writeError(w, "could not load emotes", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, emoteList)
	return
}

func (a *API) voteEmote(w http.ResponseWriter, r *http.Request) {
	// Must be POST
	if r.Method != "POST" {
		w.WriteHeader(404)
		return
	}

	// Load emote list
	if emoteList == nil {
		var err error
		emoteList, err = loadEmotes()
		if err != nil {
			writeError(w, "could not load emotes: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Load user from token
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		writeError(w, "could not load user data from request context", http.StatusInternalServerError)
		return
	}
	// Load user from database
	if err := a.LoadUser(r.Context(), user); err != nil {
		writeError(w, "could not load user data from database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse Form
	if err := r.ParseForm(); err != nil {
		writeError(w, "could not parse post form: "+err.Error(), http.StatusInternalServerError)
		return
	}

	emote := r.FormValue("emote")
	action := r.FormValue("action")

	if findEmote(emote, emoteList) == nil {
		writeError(w, "an emote with the specified name does not exist", http.StatusBadRequest)
		return
	}
	if !(action == "add" || action == "retract") {
		writeError(w, "bad action", http.StatusBadRequest)
		return
	}
	if err := a.VoteEmote(r.Context(), user, emote, action); err != nil {
		writeError(w, "could not update emote vote: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, "")
	return
}

func loadEmotes() ([]*Emote, error) {
	emotes := []*Emote{}

	err := filepath.Walk(emotesPath, func(path string, info os.FileInfo, err error) error {
		filename := strings.TrimPrefix(path, emotesPath)
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		if name != "" {
			emotes = append(emotes, &Emote{name, filename, 0})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return emotes, nil
}

func findEmote(name string, emotes []*Emote) *Emote {
	for _, emote := range emotes {
		if emote.Name == name {
			return emote
		}
	}
	return nil
}
