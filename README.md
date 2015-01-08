# CSV


## Introduction

This small package provides a CSV Encoder similar to json.Encoder or xml.Encoder
of the standard library.

It is an early version and currently only handles slices of structs. So pull
requests are welcomed to fix bugs or handle more use cases (slices or maps of
string, etc.).


## Documentation

http://godoc.org/github.com/SoCloz/csv


## Download

`go get github.com/SoCloz/csv`


## Example

```go
package main

import (
	"github.com/SoCloz/csv"
)

func main() {
	v := []struct {
		Name string
		Age  int
	}{
		{Name: "Bob", Age: 42},
		{Name: "Joe", Age: 17},
	}
	err := csv.NewEncoder(os.Stdout).Encode(v)
	if err != nil {
		panic(err)
	}
	// Output:
	// Name,Age
	// Bob,42
	// Joe,17
}
```
