package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dupakarovsky/todo"
)

//=================================
// IMPROVING THE LIST OUTPUT FORMAT
//=================================
// We can improove the output format by implementing a Stringer interaface to our List type.
// this allow us to pass our type to any formating function that expects a string

// > add a String() method to the List type in todo.go

// hardcode the filename for now
const todoFileName = "todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by Dupakarovksy\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage Information:\n")
		flag.PrintDefaults()
	}

	// Add a list of flags to be passed to the command line
	task := flag.String("task", "", "Task to be included in the ToDo list")
	list := flag.Bool("list", false, "List all ToDo items")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	// define a instance of a Todo List initialize in it's zero value
	l := &todo.List{}

	// try to read the todoFileName using the Get() method.
	if err := l.Get(todoFileName); err != nil {
		// if fails, print the error to the Standard Error in Terminal
		fmt.Fprintln(os.Stderr, err)
		// exit the process with code 1 (error condition)
		os.Exit(1)
	}

	// File doesn't exist or file was successfuly read:
	// check if any arguments were passed to the command line
	switch {
	//  check for the case where the '-list' flag is passed
	case *list:
		// INFO: now we can replace the for-range loop with a call to Print(l)
		// Pritnt(l) uses the default String() for the type, which in our case uses a for-range
		// loop to consturct a formated output.
		// ERROR: the main_test.ListTask should fail now, as we're chaning the output.
		// > go there to correct it.

		fmt.Print(l)

		// check for the case where the '-complete' flag is passed
	case *complete > 0:
		// call the Complete() to update Done and CompletedAt fields
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// save the updated list on disk.
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// check for the case where the '-task' flag is passed
	case *task != "":
		// call the Add method to add the string as a new task  name
		l.Add(*task)

		// save the updated list on disk.
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		// update the default case to output an error to stderr
	default:
		// Check for error during save.
		fmt.Fprintln(os.Stderr, "Invalid Option")
		// exit the program if errors saving
		os.Exit(1)
	}
}
