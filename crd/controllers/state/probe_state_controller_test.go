package state

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
)

var _ = Describe("Probe State", func() {
	var sm *StateManager

	BeforeEach(func() {
		// Create a fresh state manager for each test to ensure isolation
		sm = NewStateManager()
	})

	It("Records state", func() {
		initialState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "abc",
			End:                        "",
		}

		err := sm.SetProbeState(&initialState)
		Expect(err).To(BeNil())

		retrievedState := sm.GetProbeState()

		// Compare field by field for clarity
		Expect(retrievedState.Items).To(Equal(initialState.Items))
		Expect(retrievedState.Errors).To(Equal(initialState.Errors))
		Expect(retrievedState.PodNetworking).To(Equal(initialState.PodNetworking))
		Expect(retrievedState.Start).To(Equal("abc"))
		Expect(retrievedState.End).To(Equal(""))
	})

	It("Updates State", func() {
		initialState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "abc",
			End:                        "",
		}

		finalState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "def",
			End:                        "",
		}

		err := sm.SetProbeState(&initialState)
		Expect(err).To(BeNil())

		err = sm.SetProbeState(&finalState)
		Expect(err).To(BeNil())

		retrievedState := sm.GetProbeState()

		// Verify the state was updated to finalState
		Expect(retrievedState.Start).To(Equal("def"))
		Expect(retrievedState.Start).NotTo(Equal(initialState.Start))
	})

	It("Updates State when Pod networking is probed", func() {
		initialState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "abc",
			End:                        "",
		}

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

		err := sm.SetProbeState(&initialState)
		Expect(err).To(BeNil())

		err = sm.SetProbeState(&finalState)
		Expect(err).To(BeNil())

		retrievedState := sm.GetProbeState()

		// Verify the state was updated
		Expect(retrievedState.Start).To(Equal("def"))
		Expect(retrievedState.PodNetworking).To(HaveLen(1))
		Expect(retrievedState.PodNetworking[0].PodName).To(Equal("testPod"))
		Expect(retrievedState.PodNetworking[0].Netstat).To(Equal("somestring"))

		// Verify it's different from initial state
		Expect(retrievedState.Start).NotTo(Equal(initialState.Start))
	})
})

// Tests using the default global manager (for backward compatibility testing)
var _ = Describe("Probe State with default manager", func() {
	BeforeEach(func() {
		// Reset the default manager to a clean state
		cleanState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
		}
		SetProbeState(&cleanState)
	})

	It("Records state using package-level functions", func() {
		initialState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "abc",
			End:                        "",
		}

		SetProbeState(&initialState)
		retrievedState := GetProbeState()

		Expect(retrievedState.Start).To(Equal("abc"))
		Expect(retrievedState.End).To(Equal(""))
	})

	It("Updates State using package-level functions", func() {
		initialState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "abc",
			End:                        "",
		}

		finalState := v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
			Start:                      "def",
			End:                        "",
		}

		SetProbeState(&initialState)
		SetProbeState(&finalState)

		retrievedState := GetProbeState()
		Expect(retrievedState.Start).To(Equal("def"))
	})
})
