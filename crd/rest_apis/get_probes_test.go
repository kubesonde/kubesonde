package restapis

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/state"
)

func TestContinuousMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest APIs")
}

var _ = Describe("GetProbes", func() {
	var stateManager *state.StateManager

	BeforeEach(func() {
		// Create a fresh state manager for each test
		stateManager = state.NewStateManager()

		// Initialize with empty state
		innerState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
		}
		err := stateManager.SetProbeState(&innerState)
		Expect(err).To(BeNil())
	})

	It("Returns probes", func() {
		finalState := v1.ProbeOutput{
			Items:  []v1.ProbeOutputItem{},
			Errors: []v1.ProbeOutputError{},
			PodNetworking: []v1.PodNetworkingInfo{
				{
					PodName: "testPod",
					Netstat: "somestring",
				},
			},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
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
		Expect(dst.PodNetworking[0].PodName).To(Equal("testPod"))
		Expect(dst.PodNetworking[0].Netstat).To(Equal("somestring"))
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
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
		}
		state.SetProbeState(&innerState)
	})

	It("Returns probes using default manager", func() {
		finalState := v1.ProbeOutput{
			Items:  []v1.ProbeOutputItem{},
			Errors: []v1.ProbeOutputError{},
			PodNetworking: []v1.PodNetworkingInfo{
				{
					PodName: "testPod",
					Netstat: "somestring",
				},
			},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
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
		Expect(dst.PodNetworking[0].PodName).To(Equal("testPod"))
		Expect(dst.PodNetworking[0].Netstat).To(Equal("somestring"))
		Expect(dst.Start).To(Equal("def"))
		Expect(dst.End).To(BeEmpty())
	})
})
