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

var _ = Describe("GetProbes", func() {
	var stateManager *state.StateManager

	BeforeEach(func() {
		// Create a fresh state manager for each test
		stateManager = state.NewStateManager()

		// Initialize with empty state
		innerState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              make(v1.PodNetworkingInfo),
			PodConfigurationNetworking: make(v1.PodNetworkingInfo),
		}
		err := stateManager.SetProbeState(&innerState)
		Expect(err).To(BeNil())
	})

	It("Returns probes", func() {
		finalState := v1.ProbeOutput{
			Items:  []v1.ProbeOutputItem{},
			Errors: []v1.ProbeOutputError{},
			PodNetworking: v1.PodNetworkingInfo{
				"testPod": {
					{Port: "80", IP: "1.2.3.4", Protocol: "TCP"},
				},
			},
			PodConfigurationNetworking: make(v1.PodNetworkingInfo),
			Start:                      "def",
			End:                        "",
		}

		err := stateManager.SetProbeState(&finalState)
		Expect(err).To(BeNil())

		req := httptest.NewRequest("GET", "http://localhost:2709/probes", nil)
		w := httptest.NewRecorder()

		// Use the handler with our custom state manager
		handler := GetProbesHandlerWithManager(stateManager)
		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(200))

		b, err := io.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		var dst v1.ProbeOutput
		err = json.Unmarshal(b, &dst)
		Expect(err).To(BeNil())

		// Compare the relevant fields
		Expect(dst.Items).To(BeEmpty())
		Expect(dst.Errors).To(BeEmpty())
		Expect(dst.PodNetworking).To(HaveLen(1))
		Expect(dst.PodNetworking["testPod"]).To(HaveLen(1))
		Expect(dst.PodNetworking["testPod"][0].Port).To(Equal("80"))
		Expect(dst.Start).To(Equal("def"))
		Expect(dst.End).To(Equal(""))
	})
})

// If you need to use the default manager instead of a custom one,
// use this alternative approach:
var _ = Describe("GetProbes with default manager", func() {
	BeforeEach(func() {
		// Reset the default manager state
		innerState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              make(v1.PodNetworkingInfo),
			PodConfigurationNetworking: make(v1.PodNetworkingInfo),
		}
		state.SetProbeState(&innerState)
	})

	It("Returns probes using default manager", func() {
		finalState := v1.ProbeOutput{
			Items:  []v1.ProbeOutputItem{},
			Errors: []v1.ProbeOutputError{},
			PodNetworking: v1.PodNetworkingInfo{
				"testPod": {
					{Port: "80", IP: "1.2.3.4", Protocol: "TCP"},
				},
			},
			PodConfigurationNetworking: make(v1.PodNetworkingInfo),
			Start:                      "def",
			End:                        "",
		}

		state.SetProbeState(&finalState)

		req := httptest.NewRequest("GET", "http://localhost:2709/probes", nil)
		w := httptest.NewRecorder()
		handler := GetProbesHandler()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(200))

		b, err := io.ReadAll(resp.Body)
		Expect(err).To(BeNil())

		var dst v1.ProbeOutput
		err = json.Unmarshal(b, &dst)
		Expect(err).To(BeNil())

		// More explicit field-by-field comparison
		Expect(dst.Items).To(BeEmpty())
		Expect(dst.Errors).To(BeEmpty())
		Expect(dst.PodNetworking).To(HaveLen(1))
		Expect(dst.PodNetworking["testPod"]).To(HaveLen(1))
		Expect(dst.PodNetworking["testPod"][0].Port).To(Equal("80"))
		Expect(dst.Start).To(Equal("def"))
		Expect(dst.End).To(BeEmpty())
	})
})
