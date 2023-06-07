package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	eventstorage "kubesonde.io/controllers/event-storage"
	"kubesonde.io/rest_apis/types"
)

func TestMonitor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "monitor")
}

var _ = Describe("findListeningPortsNonInLoopback", func() {
	It("works", func() {
		p1 := []types.NestatInfoRequestBodyItem{{
			Type:  1,
			Laddr: []string{"1.2.3.4", "80"},
		}}
		Expect(len(findListeningPortsNonInLoopback(p1))).To(Equal(1))
		p2 := []types.NestatInfoRequestBodyItem{{
			Type:  1,
			Laddr: []string{"127.0.0.1", "80"},
		}}
		Expect(len(findListeningPortsNonInLoopback(p2))).To(Equal(0))
	})
	It("finds udp and tcp", func() {
		p1 := []types.NestatInfoRequestBodyItem{{
			Type:  1,
			Laddr: []string{"1.2.3.4", "80"},
		}}
		Expect(findListeningPortsNonInLoopback(p1)[0].Protocol).To(Equal("TCP"))
		p2 := []types.NestatInfoRequestBodyItem{{
			Type:  2,
			Laddr: []string{"1.2.3.4", "80"},
		}}
		Expect(findListeningPortsNonInLoopback(p2)[0].Protocol).To(Equal("UDP"))
	})
})

var _ = Describe("buildProbesFromMonitorContainer", func() {
	It("works", func() {
		// Given
		eventstorage.AddActivePod("anotherpod", eventstorage.CreatedPodRecord{
			Pod: v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "anotherpod", Namespace: "mynamespace"}, Status: v12.PodStatus{PodIP: "1.1.1.1"}},
		})
		eventstorage.AddActivePod("testpod", eventstorage.CreatedPodRecord{
			Pod: v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod", Namespace: "mynamespace"}, Status: v12.PodStatus{PodIP: "1.2.3.4"}},
		})
		client := fake.NewSimpleClientset()
		p := &v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod", Namespace: "mynamespace"}}
		client.CoreV1().Pods("mynamespace").Create(context.TODO(), p, metav1.CreateOptions{})
		p1 := []types.NestatInfoRequestBodyItem{{
			Type:  1,
			Laddr: []string{"1.2.3.4", "80"},
		}}
		// When/then
		Expect(len(buildProbesFromMonitorContainer(client, p1, "testpod"))).To(Equal(1))

		// cleanup
		eventstorage.DeleteActivePod("anotherpod")
		eventstorage.DeleteActivePod("testpod")
	})

})

var _ = Describe("Decode netinfo data", func() {
	It("Returns error when buffer does not exist", func() {
		buffer := new(bytes.Buffer)
		_, err := eventuallyDecodeNetinfoData(buffer)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("not found"))
	})

	It("Returns error when no new line is available", func() {
		buffer := new(bytes.Buffer)
		message := "Test"
		buffer.Write([]byte(message))
		_, err := eventuallyDecodeNetinfoData(buffer)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("not found"))
	})
	It("Returns error when data structure cannot be decoded and starts with newline", func() {
		buffer := new(bytes.Buffer)
		message := "\nTest"
		buffer.Write([]byte(message))
		_, err := eventuallyDecodeNetinfoData(buffer)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("not found"))
	})
	It("Returns error when data structure cannot be decoded and ends with newline", func() {
		buffer := new(bytes.Buffer)
		message := "Test\n"
		buffer.Write([]byte(message))
		_, err := eventuallyDecodeNetinfoData(buffer)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("could not decode"))
	})
	It("Decodes netinfo information", func() {
		netinfo := types.NestatInfoRequestBody{
			types.NestatInfoRequestBodyItem{Fd: 1,
				Family: 2,
				Type:   1,
				Laddr:  []string{"1.1.1.1", "8080"}},
		}
		data := lo.Must1(json.Marshal(netinfo))
		buffer := new(bytes.Buffer)
		message := fmt.Sprintf("%s\n", data)
		buffer.Write([]byte(message))
		payload, err := eventuallyDecodeNetinfoData(buffer)
		Expect(err).To(BeNil())
		Expect(payload).To(Equal(netinfo))
	})
})
