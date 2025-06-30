package gomoon_test

import (
	"unsafe"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/0komo/gomoon/internal/tests"

	moon "github.com/0komo/gomoon"
)

var _ = //

Describe("State", Ordered, func() {
	Context("when initializing", func() {
		It("is successful", func() {
			L := moon.NewState()
			if Expect(L).ToNot(BeNil()) {
				L.Close()
			}
		})

		It("with custom alloc fn is successful", func() {
			L := moon.NewStateWithAllocFn(func(ptr unsafe.Pointer, _, nsize uintptr) unsafe.Pointer {
				if nsize == 0 {
					tests.Free(ptr)
					return nil
				}
				return tests.Realloc(ptr, nsize)
			}, nil)
			if Expect(L).ToNot(BeNil()) {
				L.Close()
			}
		})
	})

	var L *moon.State

	BeforeAll(func() {
		L = moon.NewState()
		Expect(L).ToNot(BeNil())
	})

	AfterAll(func() {
		L.Close()
		Expect(L.IsClosed()).To(BeTrue())
	})

	AfterEach(func() {
		L.Pop(L.GetTop())
		Expect(L.GetTop()).To(BeZero())
	})

	It("can retrieve a string from stack", func() {
		testStr := "foo\x00\x00"
		L.PushString(testStr)
		Expect(L.GetTop()).To(Equal(1))
		str, ok := L.ToString(-1)
		Expect(ok).To(BeTrue())
		Expect(str).To(Equal(testStr))
	})

	// It("can call a Go function", func() {
	// 	L.PushGoFunction(func(_ *moon.State) int {
	// 		L.PushString("foo")
	// 		L.PushBool(true)
	// 		return 2
	// 	})
	// 	L.SetGlobal("foo")
	// 	L.PushGoFunction(func(_ *moon.State) int {
	// 		b := L.ToBool(-1)
	// 		Expect(b).To(BeTrue())
	// 		return 0
	// 	})
	// 	L.SetGlobal("assert")
	// 	ok := L.DoString(`
	// 		local a, b = foo()
	// 		assert(a == "foo" and b == true)
	// 	`)
	// 	Expect(ok).To(BeTrue())
	// })
})
