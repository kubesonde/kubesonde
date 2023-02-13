package monitor

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
			Pod: v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "anotherpod"}, Status: v12.PodStatus{PodIP: "1.1.1.1"}},
		})
		eventstorage.AddActivePod("testpod", eventstorage.CreatedPodRecord{
			Pod: v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod"}, Status: v12.PodStatus{PodIP: "1.2.3.4"}},
		})
		client := fake.NewSimpleClientset()
		p := &v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "testpod"}}
		client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		p1 := []types.NestatInfoRequestBodyItem{{
			Type:  1,
			Laddr: []string{"1.2.3.4", "80"},
		}}
		// When/then
		Expect(len(buildProbesFromMonitorContainer(client, p1, "testpod"))).To(Equal(1))

		//cleanup
		eventstorage.DeleteActivePod("anotherpod")
		eventstorage.DeleteActivePod("testpod")
	})

})
