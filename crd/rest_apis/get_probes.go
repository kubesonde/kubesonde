package restapis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"kubesonde.io/controllers/state"
)

const GET_PROBES_PATH = "/probes"

func GetProbesHandler() http.Handler {
	handerFun := func(w http.ResponseWriter, r *http.Request) {
		var currState = state.GetProbeState()
		data, err := json.MarshalIndent(currState, "", "  ")
		if err != nil {
			log.Error(err, "[GET /probes]Could not not return data")
			fmt.Fprintf(w, "Error")
		}
		fmt.Fprintf(w, string(data))
	}
	handler := http.HandlerFunc(handerFun)
	return handler

}
