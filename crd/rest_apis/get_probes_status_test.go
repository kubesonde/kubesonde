package restapis

import (
	"encoding/json"
	"io"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/state"
)

var _ = Describe("GetProbesStatus", func() {
	var stateManager *state.StateManager

	BeforeEach(func() {
		stateManager = state.NewStateManager()
	})

	It("rejects non-GET methods", func() {
		req := httptest.NewRequest("POST", "http://localhost:2709/probes/status", nil)
		w := httptest.NewRecorder()

		GetProbesStatusHandlerWithManager(stateManager).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(405))
	})

	It("reports not-complete before any probe is recorded", func() {
		req := httptest.NewRequest("GET", "http://localhost:2709/probes/status", nil)
		w := httptest.NewRecorder()

		GetProbesStatusHandlerWithManager(stateManager).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))

		b, err := io.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		var status state.Completeness
		Expect(json.Unmarshal(b, &status)).To(BeNil())
		Expect(status.Complete).To(BeFalse())
		Expect(status.Count).To(Equal(0))
	})

	It("reports the current probe count", func() {
		items := []v1.ProbeOutputItem{{Port: "80"}}
		Expect(stateManager.AppendProbes(&items)).To(BeNil())

		req := httptest.NewRequest("GET", "http://localhost:2709/probes/status", nil)
		w := httptest.NewRecorder()

		GetProbesStatusHandlerWithManager(stateManager).ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))

		b, err := io.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		var status state.Completeness
		Expect(json.Unmarshal(b, &status)).To(BeNil())
		Expect(status.Count).To(Equal(1))
	})
})
