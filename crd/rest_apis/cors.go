package restapis

import "net/http"

// withCORS wraps an http.Handler with permissive CORS headers so the frontend
// (served from a different origin) can call the controller REST API directly.
// It also short-circuits OPTIONS preflight requests with a 204 before the
// wrapped handler runs, so method checks in the handlers don't reject them.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
