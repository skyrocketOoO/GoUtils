package main

import (
	"fmt"
	"reflect"
)

// Struct C definition
type C struct {
	Val string
}

// Struct B definition
type B struct {
	C     *C
	Value string
}

// Struct A definition
type A struct {
	B *B
}

// NewA creates a new instance of A with initialized fields using reflection
func NewA() *A {
	a := new(A)
	InitializeFields(a)
	return a
}

// initializeFields initializes the fields of A using reflection
func InitializeFields(a any) {
	// Get the reflect.Value of A
	val := reflect.ValueOf(a).Elem()

	// Iterate through fields of A
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field is a pointer and nil
		if field.Kind() == reflect.Ptr && field.IsNil() {
			// Allocate memory for the pointer field
			field.Set(reflect.New(field.Type().Elem()))

			// If the field is a struct pointer, recursively initialize its fields
			if field.Type().Elem().Kind() == reflect.Struct {
				// Get the actual value of the pointer
				fieldVal := field.Elem()

				// Initialize fields of the nested struct recursively
				initializeStructFields(fieldVal)
			}
		}
	}
}

// initializeStructFields initializes the fields of a struct using reflection recursively
func initializeStructFields(val reflect.Value) {
	// Iterate through fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Check if the field is a pointer and nil
		if field.Kind() == reflect.Ptr && field.IsNil() {
			// Allocate memory for the pointer field
			field.Set(reflect.New(field.Type().Elem()))

			// If the field is a struct pointer, recursively initialize its fields
			if field.Type().Elem().Kind() == reflect.Struct {
				// Get the actual value of the pointer
				fieldVal := field.Elem()

				// Recursively initialize fields of the nested struct
				initializeStructFields(fieldVal)
			}
		}
	}
}

func main() {
	// Create a new instance of A
	a := NewA()

	// Print the value of a.B.C.Val
	fmt.Println(a.B.C) // This will print an empty string ""
}
