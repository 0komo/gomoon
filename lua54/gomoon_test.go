package gomoon_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLua54(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lua54 Suite")
}
