package eventstorage

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEventStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EventStorage race conditions")
}
