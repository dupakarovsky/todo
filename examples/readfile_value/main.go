package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	ConcludedAt time.Time
}

type List []item

func (l List) Get(fileName string) error {
	// read json file from disk if exists
	file, err := os.ReadFile(fileName)
	if err != nil {
		switch {
		// file doesnt exist
		case errors.Is(err, os.ErrNotExist):
			return fmt.Errorf("ErrNotExist: %+w", err)
			// unknown error
		default:
			return err
		}
	}

	fmt.Println("File successfuly read")

	i := item{
		Task:        string(file),
		Done:        false,
		CreatedAt:   time.Now(),
		ConcludedAt: time.Time{},
	}

	l = append(l, i)
	fmt.Printf("-- i:\n")
	fmt.Printf("\tvalue: [ %v ]\n", i)
	fmt.Printf("\taddr : [ %p ]\n", &i)

	return nil
}

func main() {
	// l with zero value
	var l List

	fmt.Println("-- l info ---- ")
	fmt.Printf("\tl: value: [ %v ]\n", l)
	fmt.Printf("\tl: addr : [ %p ]\n", &l)

	// call to Get()
	if err := l.Get("todo.json"); err != nil {
		fmt.Println(err)
		return
	}

	// print
	fmt.Println("Program has exited")
}
