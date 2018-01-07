package pyfmt_test

import (
	"fmt"

	"github.com/slongfield/pyfmt"
)

func ExampleMust() {
	fmt.Println(pyfmt.Must("{}", 3))
	fmt.Println(pyfmt.Must("{test}", struct{ test int }{test: 42}))
	fmt.Println(pyfmt.Must("{bar}{foo}{1}", struct {
		foo string
		bar string
	}{foo: "世界", bar: "你好"}, "!"))
	fmt.Println(pyfmt.Must(`"{2:^7}"`, 1, 12, 123, 1234))
	fmt.Println(pyfmt.Must("{:#x}", 3735928559))
	// Output:
	// 3
	// 42
	// 你好世界!
	// "  123  "
	// 0xdeadbeef
}
