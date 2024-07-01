package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
)

//=====================
// TODO API
//=====================

// 'item' will be used internally by the 'todo' package.
// will hold fields representing a information about a particular 'ToDo' item.
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represents a list of Todo items
type List []item

// Add will add a new ToDo element to the List slice
func (l *List) Add(taskName string) {

	// instantiate a new ToDo item.
	td := item{
		Task:        taskName,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{}, // zero value of the time.Time struct
	}

	// append the td into the slice.
	//INFO: a call to append needs to be done with a slice. So we're deferencenging the pointer to access the value
	*l = append(*l, td)
}

// Complete method will mark a ToDo item in the list as completed
// INFO: the method will modifying the List 'l'. So using a pointer receiver is just to maintain conistency on the methods.
func (l *List) Complete(pos int) error {

	// store the dereferenced value of the List l to perform a len check
	// the backing array of ls and l are the same.
	ls := *l

	// check whether the position passed is valid
	if pos <= 0 || pos < len(ls) {
		return fmt.Errorf("item %d does not exist", pos)
	}

	// Update the ToDo Done and CompleteAt fields
	// Backing Array is the same for ls and l. ls is modifying the backing array, l will change as well.
	ls[pos-1].Done = true
	ls[pos-1].CompletedAt = time.Now()

	return nil
}

// Delete will remove an Todo item from the List
func (l *List) Delete(pos int) error {
	// store the dereferenced value of the List l to perform a len check
	// both ls and l have the same Backing Array
	ls := *l

	// check whether the pos passed is valid
	if pos <= 0 || pos > len(ls) {
		return fmt.Errorf("item %d does not exist", pos)
	}

	// remove the ToDo from the slice by slicing off it's position (unlinke index, position starts at 1)
	left := ls[:pos-1] // first half of the slice, without the element we need to cut
	right := ls[pos:]  // second half of the slice

	// update the slice by appending the two halfs
	// append copies by value, and may create a new Backing Array. So copy the new slice to the dereferenced value of l
	*l = append(left, right...)

	return nil
}

// Save method will encode the List as JSON an save it using the provided filename
func (l *List) Save(filename string) error {
	// marshal the List into json format
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	// write to the file system
	return os.WriteFile(filename, js, fs.FileMode(os.O_WRONLY))
}

// Get method will open the file and decode the json file into the List slice
func (l *List) Get(filename string) error {

	// try read the file from the os.
	file, err := os.ReadFile(filename)
	if err != nil {
		switch {
		// file didn't exist
		case errors.Is(err, os.ErrNotExist):
			return nil
			// some other unknown error
		default:
			return err
		}
	}
	// check wether the file is empty
	if len(file) == 0 {
		return nil
	}

	// file read. // Unmarshal form JSON into the List slice.
	return json.Unmarshal(file, &l)
}
