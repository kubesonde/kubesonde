package restapis

import (
	"encoding/json"
	"io/ioutil"
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
	BeforeEach(func() {
		innerState := v1.ProbeOutput{
			Items:           []v1.ProbeOutputItem{},
			Errors:          []v1.ProbeOutputError{},
			PodNetworking:   []v1.PodNetworkingInfo{},
			PodNetworkingV2: make(v1.PodNetworkingInfoV2),
		}
		state.SetProbeState(&innerState)
	})
	It("Returns probes", func() {
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
		state.SetProbeState(&finalState)
		req := httptest.NewRequest("GET", "http://localhost:2709/probes", nil)
		w := httptest.NewRecorder()
		handler := GetProbesHandler()
		handler.ServeHTTP(w, req)
		resp := w.Result()
		Expect(resp.StatusCode).To(BeIdenticalTo(200))

		b, err := ioutil.ReadAll(resp.Body)
		Expect(err).To(BeNil())
		var dst v1.ProbeOutput
		err = json.Unmarshal(b, &dst)
		Expect(err).To(BeNil())
		Expect(dst).To(BeEquivalentTo(finalState))
	})
})
