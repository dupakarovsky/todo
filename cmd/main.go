package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dupakarovsky/todo"
)

//========================
// RUNNING THE CLI
//========================

// run the CLI: /todo/cmd
// $ go run main.go

// The first run, without arguments shouldn't produce any results.
// Add some arguments

// $ go run main.go Add this to do item to the list

// This will add generate a new JSON file with one item with this title.
// If we run again the same command, the Task name will appern in the stdout

// $ go run main.go
// Add this to do item to the list

//========================================
// HANDLING MULTIPLE COMMAND LINE OPTIONS
//========================================
/*
   Use the flag to pass multiple commands to the terminal
   -list: bool. when passed will list all to-do items.
   -task: string. when passed will include a string argument as a new ToDo
   -complete: int. when passed will mark the item number as completed.
*/

// hardcode the filename for now
const todoFileName = "todo.json"

func main() {
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

// > update the test cases in main_test.go
