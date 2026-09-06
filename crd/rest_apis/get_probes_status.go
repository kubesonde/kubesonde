package restapis

import (
	"encoding/json"
	"net/http"

	"kubesonde.io/controllers/state"
)

const GET_PROBES_STATUS_PATH = "/probes/status"

// GetProbesStatusHandler reports whether probing has quiesced. Because there is
// no reliable way to know in advance how many probes will run, completeness is
// inferred from the recorded probe count being stable over a fixed window.
func GetProbesStatusHandler() http.Handler {
	return GetProbesStatusHandlerWithManager(state.GetDefaultManager())
}

func GetProbesStatusHandlerWithManager(sm *state.StateManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := sm.Completeness()
		w.Header().Set("Content-Type", "application/json")

		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			log.Error(err, "[GET /probes/status] Failed to marshal completeness")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			log.Error(err, "[GET /probes/status] Failed to write response")
		}
	})
}
