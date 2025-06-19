package gomoon_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gomoon "github.com/0komo/gomoon/lua54"
)

func clearStack(L *gomoon.State) {
	L.Pop(L.GetTop())
}

var _ = Describe("State", Ordered, func() {
	var L *gomoon.State

	BeforeAll(func() {
		raw, err := gomoon.NewState()
		L = &raw

		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		clearStack(L)
	})

	It("initialized correctly", func() {
		Expect(L.IsClosed()).To(BeFalse())
	})

	Context("can push a", func() {
		It("nil", func() {
			L.PushNil()
			typ := L.Type(-1)

			Expect(L.GetTop()).To(Equal(1))
			Expect(typ).To(Equal(gomoon.Nil))
			Expect(typ.String()).To(Equal("nil"))
		})
	})

	AfterAll(func() {
		L.Close()

		Expect(L.IsClosed()).To(BeTrue())
	})
})
