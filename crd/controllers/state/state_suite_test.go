package state

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProbeStateController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProbeState controller")
}
