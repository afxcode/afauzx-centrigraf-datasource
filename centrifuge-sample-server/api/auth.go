package api

import (
	"net/http"

	"github.com/centrifugal/centrifuge"
)

// Authentication middleware example. Centrifuge expects Credentials
// with current user ID set. Without provided Credentials client
// connection won't be accepted. Another way to authenticate connection
// is reacting to node.OnConnecting event where you may authenticate
// connection based on a custom token sent by a client in first protocol
// frame. See _examples folder in repo to find real-life auth samples
// (OAuth2, Gin sessions, JWT etc).
func auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Put authentication Credentials into request Context.
		// Since we don't have any session backend here we simply
		// set user ID as empty string. Users with empty ID called
		// anonymous users, in real app you should decide whether
		// anonymous users allowed to connect to your server or not.
		cred := &centrifuge.Credentials{
			UserID: "",
		}
		newCtx := centrifuge.SetCredentials(ctx, cred)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}
