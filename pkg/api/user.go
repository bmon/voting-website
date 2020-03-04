package api

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	Sub           string   `firestore:"sub"`
	Aud           string   `firestore:"-"`
	Azp           string   `firestore:"-"`
	Email         string   `firestore:"email"`
	EmailVerified bool     `firestore:"-"`
	Exp           int64    `firestore:"-"`
	FamilyName    string   `firestore:"-"`
	GivenName     string   `firestore:"-"`
	Hd            string   `firestore:"-"`
	Iat           int64    `firestore:"-"`
	Jti           string   `firestore:"-"`
	Locale        string   `firestore:"-"`
	Name          string   `firestore:"name"`
	Picture       string   `firestore:"-"`
	Votes         []string `firestore:"votes"`
	Admin         bool     `firestore:"-"`
}

func (a *API) LoadUser(ctx context.Context, u *User) error {
	snap, err := a.store.Collection("Users").Doc(u.Sub).Get(ctx)
	if err != nil {
		// ignore not-found errors
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	if err := snap.DataTo(u); err != nil {
		return err
	}
	return nil
}

func (a *API) listVotes(w http.ResponseWriter, r *http.Request) {
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

	// ensure the array comes out as empty array and not Null in json
	if user.Votes == nil {
		user.Votes = make([]string, 0)
	}
	writeJSON(w, user.Votes)
}

func (a *API) VoteEmote(ctx context.Context, u *User, emote string, action string) error {
	var delta int
	var arrayChange interface{}
	switch action {
	case "add":
		delta = 1
		arrayChange = firestore.ArrayUnion(emote)
	case "retract":
		delta = -1
		arrayChange = firestore.ArrayRemove(emote)
	default:
		return fmt.Errorf("VoteEmote: bad action supplied")
	}

	return a.store.RunTransaction(ctx, func(context.Context, *firestore.Transaction) error {
		ref := a.store.Collection("Users").Doc(u.Sub)
		snap, err := ref.Get(ctx)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		result, err := ref.Set(ctx, map[string]interface{}{
			"sub":   u.Sub,
			"name":  u.Name,
			"email": u.Email,
			"votes": arrayChange,
		}, firestore.MergeAll)
		if err != nil {
			return err
		}

		// if there was an update (i.e the user's votes actually changed), update the emote count.
		if result.UpdateTime.After(snap.UpdateTime) {
			_, err := a.store.Collection("Emotes").Doc(emote).Set(ctx, map[string]interface{}{
				"name":  emote,
				"count": firestore.Increment(delta),
			}, firestore.MergeAll)
			if err != nil {
				return err
			}
		}

		return nil
	})

}
