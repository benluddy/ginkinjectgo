package ginkinjectgo

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestExample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "example")
}

type foo string

type bar struct {
	a, b, c int
}

var _ = BeforeSuite(func() {
	RegisterProvider(func() foo {
		return foo("hello world")
	})
	RegisterProvider(func() *bar {
		return &bar{
			a: 1,
			b: 3,
			c: 5,
		}
	})
	RegisterProvider("baz")
})

var _ = Describe("how this works", func() {
	It("should be called with parameters injected", func(f foo, b *bar, z string) {
		Expect(f).To(Equal(foo("hello world")))
		Expect(b).NotTo(BeNil())
		Expect(*b).To(Equal(bar{a: 1, b: 3, c: 5}))
		Expect(z).To(Equal("baz"))
	})

	Describe("nested contexts", func() {
		BeforeEach(func() {
			RegisterProvider(false)
			RegisterProvider((*bar)(nil))
		})

		It("should be called with parameters injected", func(f foo, b *bar, z string, t bool) {
			Expect(f).To(Equal(foo("hello world")))
			Expect(b).To(BeNil())
			Expect(z).To(Equal("baz"))
			Expect(t).To(BeFalse())
		})
	})
})
