package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("FilterPodsByStatus", func() {
	It("Filters pod by status", func() {
		podList := v1.PodList{
			Items: []v1.Pod{
				{
					Status: v1.PodStatus{
						Phase: "target",
					},
				},
				{
					Status: v1.PodStatus{
						Phase: "target",
					},
				},
				{
					Status: v1.PodStatus{
						Phase: "another",
					},
				},
			},
		}
		result := FilterPodsByStatus(&podList, "target")
		Expect(len(result.Items)).To(Equal(2))
	})
})

var _ = Describe("Contains", func() {
	It("Returns true if contained", func() {
		Expect(contains([]string{"test", "me"}, "test")).To(BeTrue())
	})
	It("Returns false if not contained", func() {
		Expect(contains([]string{"test", "me"}, "foo")).To(BeFalse())
	})
})

var _ = Describe("InNamespace", func() {
	It("Returns true if keywords are used", func() {
		Expect(InNamespace("", "bla bla")).To(BeTrue())
		Expect(InNamespace("all", "bla bla")).To(BeTrue())
	})
	It("Returns true if same namespace", func() {
		Expect(InNamespace("bla", "bla")).To(BeTrue())
	})
	It("Returns false otherwise", func() {
		Expect(InNamespace("bla", "foo")).To(BeFalse())
	})
})

var _ = Describe("GetDeploymentNamesInNamespace", func() {
	It("Returns empty values if no deployments", func() {
		client := fake.NewSimpleClientset()
		Expect(GetDeploymentNamesInNamespace(client, "default")).To(Equal([]string{}))
	})
})
var _ = Describe("GetReplicaSetsNamesInNamespace", func() {
	It("Returns empty values if no replicas", func() {
		client := fake.NewSimpleClientset()
		Expect(GetReplicaSetsNamesInNamespace(client, "default")).To(Equal([]string{}))
	})
})
