package gomoon_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGomoon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomoon Suite")
}
