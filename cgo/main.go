package main

/*
#include "mylib.c"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	passGoStringToC()
	passStructToC()
}

func passStructToC() {
	fmt.Println("Get and pass struct")
	street := "Sezamkowa"
	city := "Warszawa"
	firstName := "Jan"
	lastName := "Kowalski"
	age := 25
	cStreet := C.CString(street)
	cCity := C.CString(city)
	cFirstName := C.CString(firstName)
	cLastName := C.CString(lastName)
	cAge := C.int(age)

	address := (*C.struct_address)(C.malloc(C.sizeof_struct_address))
	address.street = cStreet
	address.city = cCity
	person := C.struct_person{}
	person.firstName = cFirstName
	person.lastName = cLastName
	person.age = cAge
	person.address = address
	C.print_person(&person)

	C.free(unsafe.Pointer(cStreet))
	C.free(unsafe.Pointer(cCity))
	C.free(unsafe.Pointer(cFirstName))
	C.free(unsafe.Pointer(cLastName))
	C.free(unsafe.Pointer(address))
}

func passGoStringToC() {
	fmt.Println("Pass string to C function")
	str := "lorem ipsum"
	cStr := C.CString(str)
	C.print_string(cStr)
	C.free(unsafe.Pointer(cStr))
	fmt.Println()
}
