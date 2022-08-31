package state

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
)

func TestProbeStateController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProbeState")
}

var _ = Describe("Probe State", func() {
	It("Records state", func() {
		initialState := v1.ProbeOutput{
			Items:  nil,
			Errors: nil,
			Start:  "abc",
			End:    "",
		}
		SetProbeState(&initialState)
		Expect(initialState).To(Equal(GetProbeState()))
	})
	It("Updates State", func() {
		initialState := v1.ProbeOutput{
			Items:  nil,
			Errors: nil,
			Start:  "abc",
			End:    "",
		}
		finalState := v1.ProbeOutput{
			Items:  nil,
			Errors: nil,
			Start:  "def",
			End:    "",
		}
		SetProbeState(&initialState)
		SetProbeState(&finalState)
		Expect(initialState).NotTo(Equal(GetProbeState()))
		Expect(finalState).To(Equal(GetProbeState()))
	})
	It("Updates State when Pod networking is probed", func() {
		initialState := v1.ProbeOutput{
			Items:  nil,
			Errors: nil,
			Start:  "abc",
			End:    "",
		}
		finalState := v1.ProbeOutput{
			Items:  nil,
			Errors: nil,
			PodNetworking: []v1.PodNetworkingInfo{
				{
					PodName: "testPod",
					Netstat: "somestring",
				},
			},
			Start: "def",
			End:   "",
		}
		SetProbeState(&initialState)
		SetProbeState(&finalState)
		Expect(initialState).NotTo(Equal(GetProbeState()))
		Expect(finalState).To(Equal(GetProbeState()))
	})
	It("AppendNetInfo", func() {
		initialState := v1.ProbeOutput{
			Items:  nil,
			Errors: nil,
			Start:  "abc",
			End:    "",
		}
		netinfo := []v1.PodNetworkingInfo{
			{
				PodName: "testPod",
				Netstat: "somestring",
			},
			{
				PodName: "testPod2",
				Netstat: "somestring2",
			},
		}
		SetProbeState(&initialState)
		AppendNetInfo(&netinfo)
		Expect(netinfo).To(Equal(GetProbeState().PodNetworking))

	})

})
