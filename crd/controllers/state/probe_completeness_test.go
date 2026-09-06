package state

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "kubesonde.io/api/v1"
)

var _ = Describe("Completeness", func() {
	var sm *StateManager

	BeforeEach(func() {
		sm = NewStateManager()
	})

	It("is not complete before any probe is recorded", func() {
		c := sm.CompletenessWithin(30 * time.Second)
		Expect(c.Complete).To(BeFalse())
		Expect(c.Count).To(Equal(0))
		Expect(c.SecondsSinceLastChange).To(Equal(float64(-1)))
	})

	It("is not complete right after a probe is recorded", func() {
		items := []v1.ProbeOutputItem{{Port: "80"}}
		Expect(sm.AppendProbes(&items)).To(BeNil())

		c := sm.CompletenessWithin(30 * time.Second)
		Expect(c.Complete).To(BeFalse())
		Expect(c.Count).To(Equal(1))
		Expect(c.SecondsSinceLastChange).To(BeNumerically(">=", 0))
	})

	It("is complete once the count is stable for the window", func() {
		items := []v1.ProbeOutputItem{{Port: "80"}}
		Expect(sm.AppendProbes(&items)).To(BeNil())

		// A tiny window makes the count look stable almost immediately, which is
		// the same signal a 30s window gives after 30s of quiescence.
		Eventually(func() bool {
			return sm.CompletenessWithin(1 * time.Millisecond).Complete
		}, 2*time.Second, 10*time.Millisecond).Should(BeTrue())
	})

	It("resets to incomplete when a new probe arrives", func() {
		first := []v1.ProbeOutputItem{{Port: "80"}}
		Expect(sm.AppendProbes(&first)).To(BeNil())
		Eventually(func() bool {
			return sm.CompletenessWithin(1 * time.Millisecond).Complete
		}, 2*time.Second, 10*time.Millisecond).Should(BeTrue())

		second := []v1.ProbeOutputItem{{Port: "443"}}
		Expect(sm.AppendProbes(&second)).To(BeNil())
		Expect(sm.CompletenessWithin(30 * time.Second).Complete).To(BeFalse())
	})

	It("ignores duplicate probes that do not change the count", func() {
		items := []v1.ProbeOutputItem{{Port: "80"}}
		Expect(sm.AppendProbes(&items)).To(BeNil())
		Eventually(func() bool {
			return sm.CompletenessWithin(1 * time.Millisecond).Complete
		}, 2*time.Second, 10*time.Millisecond).Should(BeTrue())

		// Re-appending the identical probe is de-duplicated, so the count does
		// not grow and completeness must not reset.
		dup := []v1.ProbeOutputItem{{Port: "80"}}
		Expect(sm.AppendProbes(&dup)).To(BeNil())
		Expect(sm.CompletenessWithin(1 * time.Millisecond).Complete).To(BeTrue())
	})
})
