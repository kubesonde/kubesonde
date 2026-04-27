package restapis

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/state"
)

var _ = Describe("PostProbesClear", func() {
	var stateManager *state.StateManager

	BeforeEach(func() {
		stateManager = state.NewStateManager()
	})

	AfterEach(func() {
		stateManager.ClearState()
	})

	BeforeEach(func() {
		innerState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
		}
		stateManager.SetProbeState(&innerState)
	})

	It("Clears all probe state on POST", func() {
		innerState := v1.ProbeOutput{
			Items: []v1.ProbeOutputItem{
				{Type: "Probe", Source: v1.ProbeEndpointInfo{Name: "src"}},
			},
			Errors: []v1.ProbeOutputError{
				{Value: v1.ProbeOutputItem{Type: "Probe"}},
			},
			PodNetworking: []v1.PodNetworkingInfo{
				{PodName: "testPod", Netstat: "data"},
			},
			PodNetworkingV2: v1.PodNetworkingInfoV2{
				"key1": []v1.PodNetworkingItem{{Port: "80", IP: "10.0.0.1"}},
			},
			PodConfigurationNetworking: v1.PodNetworkingInfoV2{
				"pod1": []v1.PodNetworkingItem{{Port: "443", IP: "10.0.0.2"}},
			},
			Start: "startTime",
			End:   "endTime",
		}
		stateManager.SetProbeState(&innerState)

		req := httptest.NewRequest("POST", "http://localhost:2709/probes/clear", nil)
		w := httptest.NewRecorder()
		handler := PostProbesClearHandlerWithManager(stateManager)
		handler.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(200))

		result := stateManager.GetProbeState()
		Expect(result.Items).To(BeEmpty())
		Expect(result.Errors).To(BeEmpty())
		Expect(result.PodNetworking).To(BeEmpty())
		Expect(result.PodNetworkingV2).To(BeEmpty())
		Expect(result.PodConfigurationNetworking).To(BeEmpty())
		Expect(result.Start).To(BeEmpty())
		Expect(result.End).To(BeEmpty())
	})

	It("Returns 405 for non-POST methods", func() {
		req := httptest.NewRequest("GET", "http://localhost:2709/probes/clear", nil)
		w := httptest.NewRecorder()
		handler := PostProbesClearHandler()
		handler.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(405))
	})
})

var _ = Describe("PostProbesClear with default manager", func() {
	BeforeEach(func() {
		state.ResetDefaultManager()
		innerState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
		}
		state.SetProbeState(&innerState)
	})

	It("Clears default manager state", func() {
		innerState := v1.ProbeOutput{
			Items: []v1.ProbeOutputItem{
				{Type: "Probe", Source: v1.ProbeEndpointInfo{Name: "src"}},
			},
			Start: "startTime",
			End:   "endTime",
		}
		state.SetProbeState(&innerState)

		state.SetNestatPod("testPod")

		result := state.GetProbeState()
		Expect(result.Items).ToNot(BeEmpty())
		Expect(result.Start).To(Equal("startTime"))

		req := httptest.NewRequest("POST", "http://localhost:2709/probes/clear", nil)
		w := httptest.NewRecorder()
		handler := PostProbesClearHandler()
		handler.ServeHTTP(w, req)

		Expect(w.Code).To(Equal(200))

		result = state.GetProbeState()
		Expect(result.Items).To(BeEmpty())
		Expect(result.Start).To(BeEmpty())

		Expect(state.GetNetstatPods()).To(BeEmpty())
	})
})
