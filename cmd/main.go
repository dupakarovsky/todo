package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dupakarovsky/todo"
)

//=============================
// DISPLAY COMMMAND LINE USAGE
//=============================
/*
   The flag package provides automatic information with the -h options

   > Run the build command and than run the binary with the -h flag:
   $ go build -o=./bin/todo ./cmd
   $ cd /bin
   $./todo -h

    it also provides helpul output for invalid flags.
    $ ./todo -test  // inform this flag is inavlid.

    flag.Usage
    ----------
    It's a variable pointing to a function. We can use it to display custom information

*/

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
		for _, item := range *l {
			// list all todo elements in the List which are NOT Done
			if !item.Done {
				fmt.Println(item.Task)
			}
		}
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
