package main

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"reflect"
	"strconv"
)

// @MyAnnotation This text is description of User struct
type User struct {
	fullName string `mytag:"foo"`
	age      int    `myIntTag:"123"`
}

// @MyAnnotation This text is description of main function
func main() {
	t := reflect.TypeOf(User{})
	displayMyTagValue(t)
	displayMyIntTagValue(t)

	fset := token.NewFileSet()
	astFile, err := goparser.ParseFile(fset, "main.go", nil, goparser.ParseComments)

	if err != nil {
		panic(err)
	}

	for _, v := range astFile.Scope.Objects {
		if v.Name == "main" {
			funcDec, ok := v.Decl.(*ast.FuncDecl)
			if !ok {
				panic("The Decl is not FuncDecl")
			}

			for _, c := range funcDec.Doc.List {
				fmt.Printf("%v\n", c.Text)
			}
		}
	}
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
