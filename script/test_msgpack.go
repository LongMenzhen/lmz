package main

import (
	"github.com/cyrnicolase/lmz/model"

	msgpack "github.com/vmihailenco/msgpack/v4"
)

func main() {
	type Class struct {
		Room      int
		ClassMate int
	}

	type Foo struct {
		Name  string
		Age   int
		Class Class
	}

	foo := Foo{
		Name: "名字",
		Age:  14,
		Class: Class{
			Room:      3,
			ClassMate: 4,
		},
	}

	rds := model.Redis()

	by, _ := msgpack.Marshal(&foo)
	rds.Set("foo", by, 0)
}
