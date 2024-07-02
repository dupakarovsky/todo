package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dupakarovsky/todo"
)

//========================
// TODO COMMAND LINE TOOL
//========================
// The command line tool will be wrapper round the API.
/*
   Features
   ---------
   when executed without arguments, will list the available ToDo items by name (itam.Task)
   when executed with one or more argument, the command will concatenate the args as a new item and add it to the list
*/

// hardcode the filename for now
const todoFileName = "todo.json"

func main() {
	fmt.Println("--main: Staring main...")

	// define a instance of a Todo List initialize in it's zero value
	l := &todo.List{}

	// try to read the todoFileName using the Get() method.
	fmt.Println("--main: Performing Get() on todo.json...")
	if err := l.Get(todoFileName); err != nil {
		// if fails, print the error to the Standard Error in Terminal
		fmt.Fprintln(os.Stderr, err)
		// exit the process with code 1 (error condition)
		os.Exit(1)
	}

	// File doesn't exist or file was successfuly read:
	// check if any arguments were passed to the command line
	switch {
	//  only the program's name was entered in the command line (no arguments).
	case len(os.Args) == 1:
		fmt.Println("--main switch: No args. Printing Tasks...")
		// list all todo elements in the List
		for _, item := range *l {
			fmt.Println(item.Task)
		}

	// one or more args where passed to the command line. (TestMain)
	default:
		fmt.Println("--main switch: Several args. Performing Join...")
		// concatenate all the command line args (expect the 1st one - the program name) into a single string
		// using the space character as separator
		tn := strings.Join(os.Args[1:], " ")

		fmt.Println("--main switch: Join completed. Performing Add()...")
		// call Add with the taskName from the concatenated string
		l.Add(tn)

		// try calling the Save() to encode to JSON and write to the file system.
		fmt.Println("--main switch: Add completed. Performing Save()...")
		if err := l.Save(todoFileName); err != nil {
			// Check for error during save.
			fmt.Fprintln(os.Stderr, err)
			// exit the program if errors saving
			os.Exit(1)
		}
	}
	fmt.Println("--Exiting Program...")
}
