//INFO: in general all files in the same directory goes into the same package. The exception is when writting tests.
// we can define a differnt package for tests to access only the exported types, variables and functions from the package you are testing.

package todo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dupakarovsky/todo"
)

// TestAdd will perform a test by adding a new 'item' entry to the 'List' slice
func TestAdd(t *testing.T) {
	fmt.Println("Starting TestAdd ...")

	// declare a ToDo List in it's zero value
	var l todo.List
	fmt.Printf("l [initial]: %+v\n", l)

	// perform a call to Add passing the task name
	taskName := "New Task"
	l.Add(taskName)
	fmt.Printf("l [after Add]: %+v\n", l)

	exp := taskName
	got := l[0].Task

	// perform a check to see if the Add method was successful. Should have a task added.
	if exp != got {
		t.Errorf("expected %q, got %q instead", exp, got)
	}

}

// TestComplete will perform a test to check wether the Done and CompletedAt properties of a task item have been
// modified. It'll start by making a call to Add() to add a new entry to the 'List' and then calling Completed() to
// modify the Done and completedAt fields.
func TestComplete(t *testing.T) {
	fmt.Println("Starting TestComplete ...")

	// declare a ToDo List in it's zero value
	var l todo.List

	// perform a call to Add passing the task name
	taskName := "New Task"
	l.Add(taskName)

	// perform a check to see if the Add method was successful. Should have a new item entry on the List 'l'
	if taskName != l[0].Task {
		t.Errorf("expected %q, got %q instead", taskName, l[0].Task)
	}

	// Examine the Done propery. Should be false at this point
	if l[0].Done {
		t.Errorf("Task should not be done at this point")
	}

	// Call the Complete method passing the position to update. (pos 1 = idx 0)
	l.Complete(1)

	// Examine the Done property. Should be true at this point
	if !l[0].Done {
		t.Errorf("Taks should be completed no")
	}
}

func TestDelete(t *testing.T) {
	fmt.Println("Starting TestDelete...")

	// initialize a new List with it's Zero value
	var l todo.List

	// define a new slice of taskNames to be added with the Add method
	tasks := []string{"Task 1", "Task 2", "Task 3"}

	// range of the slice using Value Samantics passing each each element to be added
	for _, task := range tasks {
		// call Add() in each iteration
		l.Add(task)
	}

	// perform a check wether the add was successful. Should have added 3 items
	if l[0].Task != tasks[0] {
		t.Errorf("expected %q; got %q instead", l[0].Task, tasks[0])
	}

	// Call the delete method. Should delete position 2 (index = 1)
	l.Delete(2)

	// perform a check on the slice length. Should have 2 elements now
	if len(l) != 2 {
		t.Errorf("Expected list length %d; got %d instead", 2, len(l))
	}

	// perform a check on the Task name. index 1 of l should be Task 3
	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q; got %q instead", tasks[2], l[1].Task)
	}
}

// TestSaveGet will perform a test by creating two List values, calling the Add(), Save() and Get() methods.
// It'll compare the elements of each list which should be the same at the end of the test.
func TestSaveGet(t *testing.T) {
	// declared two Lists in their Zero value
	var l1 todo.List
	var l2 todo.List

	// Add a new task to l1
	taskName := "New Task"
	l1.Add(taskName)

	// perform a check to see if the Add method was successful. Should have added one item.
	if taskName != l1[0].Task {
		t.Errorf("expected %q, got %q instead", taskName, l1[0].Task)
	}

	// create a temporary file in /tmp (empty string)
	temp, err := os.CreateTemp("", "tempfile_")
	if err != nil {
		t.Fatalf("Error creating temp file : %s", err.Error())
	}
	// remove the file form the OS when the test finishes, by pasing the file name
	defer os.Remove(temp.Name())

	// call Save with the temporary file's name and perform the error check
	if err := l1.Save(temp.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err.Error())
	}

	// call Get and try to retrieve the file and unmarshal it into the l2 List
	if err = l2.Get(temp.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err.Error())
	}

	// file save and retrieved.
	// compare both lists first elemetn.Task value. Should be the same
	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q shoud match %q task", l1[0].Task, l2[0].Task)
	}

}
