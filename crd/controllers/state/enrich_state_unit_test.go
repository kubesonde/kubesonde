package state

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
)

var _ = Describe("Enrich state", func() {
	It("Enriches a state", func() {

		initialState := v1.ProbeOutput{
			Start: "start",
			End:   "end",
			Items: []v1.ProbeOutputItem{
				{
					Source:      v1.ProbeEndpointInfo{Name: "source-123-456"},
					Destination: v1.ProbeEndpointInfo{Name: "dest-123"},
				},
			},
		}

		replicas := []string{"source-123", "dest"}
		deployments := []string{"source"}

		expectedState := v1.ProbeOutput{
			Start: "start",
			End:   "end",
			Items: []v1.ProbeOutputItem{
				{
					Source:      v1.ProbeEndpointInfo{Name: "source-123-456", ReplicaSetName: "source-123", DeploymentName: "source"},
					Destination: v1.ProbeEndpointInfo{Name: "dest-123", ReplicaSetName: "dest"},
				},
			},
		}

		Expect(EnrichState(&initialState, replicas, deployments)).To(BeEquivalentTo(&expectedState))
	})
})
