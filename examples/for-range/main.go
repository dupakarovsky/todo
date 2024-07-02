package main

import (
	"fmt"
)

type data struct {
	id  int
	val string
}

type List []data

func (l *List) Get() error {
	el := data{
		id:  1,
		val: "one",
	}

	*l = append(*l, el)
	return nil
}

func main() {

	var l List
	// l := List{}

	fmt.Println("-- l ")
	fmt.Printf("\tvalue: [ %#v ]\n", l)
	fmt.Printf("\taddr : [ %p ]\n", &l)

	_ = l.Get()

	fmt.Println("-- for range -------------")
	for idx, val := range l {
		fmt.Printf("[%d]: %v\n", idx, val)
	}

	fmt.Println("-- l after get")
	fmt.Printf("\tvalue: [ %v ]\n", l)
	fmt.Printf("\taddr : [ %p ]\n", &l)
	fmt.Println("---------------------")

}
