package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type User struct {
	fullName string `mytag:"foo"`
	age      int    `myIntTag:"123"`
}

func main() {
	t := reflect.TypeOf(User{})
	displayMyTagValue(t)
	displayMyIntTagValue(t)
}

func displayMyTagValue(t reflect.Type) {
	field, ok := t.FieldByName("fullName")

	if ok {
		fmt.Printf("The value of struct tag mytag is %v\n", field.Tag.Get("mytag"))
	}
}

func displayMyIntTagValue(t reflect.Type) {
	field, ok := t.FieldByName("age")

	if ok {
		age, err := strconv.Atoi(field.Tag.Get("myIntTag"))
		if nil != err {
			panic(err)
		}
		fmt.Printf("The value of struct tag myIntTag is %v\n", age)
	}
}
