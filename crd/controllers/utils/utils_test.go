package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	kubesondev1 "kubesonde.io/api/v1"
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

var _ = Describe("GetDeployment", func() {
	It("Returns Deployment", func() {

		replicaSet := appsv1.ReplicaSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "myreplica-abcdefg",
				OwnerReferences: []metav1.OwnerReference{
					{
						Name: "myreplica",
						Kind: "Deployment",
					},
				},
			},
		}

		// When
		depname, err := GetDeployment(replicaSet)
		// Then
		Expect(err).To(BeNil())
		Expect(depname).To(Equal("myreplica"))

	})
})

var _ = Describe("GetReplicaSet", func() {
	It("Returns ReplicaSet", func() {

		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "myreplica-abcdefg-123123",
				OwnerReferences: []metav1.OwnerReference{
					{
						Name: "myreplica-abcdefg",
						Kind: "ReplicaSet",
					},
				},
			},
		}

		// When
		depname, err := GetReplicaSet(pod)
		// Then
		Expect(err).To(BeNil())
		Expect(depname).To(Equal("myreplica-abcdefg"))

	})
})

var _ = Describe("SourcePodMatchesKubesondeSpec", func() {
	It("Returns True when Pod from different namespace is included", func() {
		ksonde := kubesondev1.Kubesonde{
			Spec: kubesondev1.KubesondeSpec{
				Namespace: "test",
				Probe:     "all",
				Include: []kubesondev1.IncludedItem{
					{
						FromPodSelector: "app=frontend",
					},
				},
			},
		}
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "another",
				Labels:    map[string]string{"app": "frontend"},
			},
		}
		Expect(SourcePodMatchesKubesondeSpec(ksonde, pod)).To(BeTrue())
	})
	It("Returns True for a Pod from the same namespace if probe all", func() {
		ksonde := kubesondev1.Kubesonde{
			Spec: kubesondev1.KubesondeSpec{
				Namespace: "test",
				Probe:     "all",
			},
		}
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Labels:    map[string]string{"app": "frontend"},
			},
		}
		Expect(SourcePodMatchesKubesondeSpec(ksonde, pod)).To(BeTrue())
	})
	It("Returns False for a Pod from the same namespace if probe none", func() {
		ksonde := kubesondev1.Kubesonde{
			Spec: kubesondev1.KubesondeSpec{
				Namespace: "test",
				Probe:     "none",
			},
		}
		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Labels:    map[string]string{"app": "frontend"},
			},
		}
		Expect(SourcePodMatchesKubesondeSpec(ksonde, pod)).To(BeFalse())
	})
})
