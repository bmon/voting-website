package api

import (
	"context"
	"fmt"
	"net/http"

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
}

func (a *API) CreateOrUpdateUser(ctx context.Context, u *User) error {
	ref := a.store.Collection("Users").Doc(u.Sub)
	_, err := ref.Set(ctx, u)
	return err
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

func (a *API) VoteEmote(ctx context.Context, u *User, emote, action string) error {
	switch action {
	case "add":
		for _, v := range u.Votes {
			if v == emote {
				// if the user already voted for it just return success - don't add it
				return nil
			}
		}
		u.Votes = append(u.Votes, emote)
		// update count
	case "retract":
		for i, v := range u.Votes {
			if v == emote {
				// copy remaining values in slice backward one element
				copy(u.Votes[i:], u.Votes[i+1:])
				// truncate slice
				u.Votes = u.Votes[:len(u.Votes)-1]
				// update count
				break
			}
		}
	default:
		return fmt.Errorf("VoteEmote: bad action supplied")
	}
	a.CreateOrUpdateUser(ctx, u)

	return nil
}
