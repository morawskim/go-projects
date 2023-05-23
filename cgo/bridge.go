package main

import "C"
import (
	"fmt"
	"strings"
)

//export screaming
func screaming(str *C.char) {
	fmt.Println(strings.ToUpper(C.GoString(str)))
}
