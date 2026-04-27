package restapis

import (
	"net/http"

	"kubesonde.io/controllers/state"
)

const POST_PROBES_CLEAR_PATH = "/probes/clear"

func PostProbesClearHandler() http.Handler {
	return PostProbesClearHandlerWithManager(state.GetDefaultManager())
}

func PostProbesClearHandlerWithManager(sm *state.StateManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sm.Clear()
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Error(err, "[POST /probes/clear] Failed to write response")
		}
	})
}
