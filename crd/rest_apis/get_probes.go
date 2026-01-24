package restapis

import (
	"encoding/json"
	"net/http"

	"kubesonde.io/controllers/state"
)

const GET_PROBES_PATH = "/probes"

func GetProbesHandler() http.Handler {
	return GetProbesHandlerWithManager(state.GetDefaultManager())
}

func GetProbesHandlerWithManager(sm *state.StateManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currState := sm.GetProbeState()
		w.Header().Set("Content-Type", "application/json")

		// Marshal state to JSON with indentation
		data, err := json.MarshalIndent(currState, "", "  ")
		if err != nil {
			log.Error(err, "[GET /probes] Failed to marshal probe state")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			log.Error(err, "[GET /probes] Failed to write response")
		}
	})
}
