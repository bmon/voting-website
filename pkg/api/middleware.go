package api

import (
	"context"
	"log"
	"net/http"
)

type User struct {
	Aud           string
	Azp           string
	Email         string
	EmailVerified bool
	Exp           int64
	FamilyName    string
	GivenName     string
	Hd            string
	Iat           int64
	Jti           string
	Locale        string
	Name          string
	Picture       string
	Sub           string
}

// logRequest middleware will log requests before handling them
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (a *API) mustVerify(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idToken, err := a.verifier.Verify(r.Context(), r.Header.Get("authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		u := &User{}
		if err = idToken.Claims(u); err != nil {
			log.Println(err)
			writeError(w, "unable to parse jwt claims", http.StatusInternalServerError)
			return
		}

		if u.Hd != a.Config.LoginDomain {
			writeError(w, "bad login domain", http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", u)))
	}
}
