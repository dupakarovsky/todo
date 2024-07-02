package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dupakarovsky/todo"
)

//=================================
// EXERCISES
//=================================
// > Implement a flag -del to delete an item from the list.
// > Add another flag to add verbose output, showing information like date/time
// > Add a flag to prevent displaying completed task.
// > Update the custom usage function to include instructions on how to provide new tasks
// > Add test cases for the other options (-complete, -delete).

// name of the json file that'll be created.
var todoFileName = "todo.json"

// getTask will accept a first parameter that implements the io.Reader interface. Then a variadict string parameter to collect all
// others arguments passd in into a slice.
func getTask(r io.Reader, args ...string) (string, error) {

	// checks whether we passed any arguments
	if len(args) > 0 {
		// arguments are provided. concatenate them into a string separeted by space and return it
		return strings.Join(args, " "), nil
	}

	// no arguments are provided. Start scanning the stdin input
	// instantiate a new Scanner to read data from the r (stdin)
	scanner := bufio.NewScanner(r)

	// scan a single line
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// check the length for the scanned line
	if len(scanner.Text()) == 0 {
		return "", fmt.Errorf("Task cannot be blank")
	}
	// return the string of the scanned line
	return scanner.Text(), nil
}

func main() {

	// the output below will be displayed when the ./todo -h is invoked.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by Dupakarovksy\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage Information:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "To add a new task use the -add flag followed by the task's name:\n(e.g: ./todo -add My New Task)\n\n")
		flag.PrintDefaults()
	}

	// Add a list of flags to be passed to the command line
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all ToDo items")
	complete := flag.Int("complete", 0, "Mark item as completed")
	// INFO: add falgs: -del, -verbose, -active
	del := flag.Int("del", 0, "Delete a task from the ToDo list")
	verbose := flag.Bool("verbose", false, "Display verbose output")
	active := flag.Bool("active", false, "Display active tasks only")

	flag.Parse()

	// check whether we have a environmental variable set. If so, set it as the value for the todoFileName var.
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

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

	// INFO: check case where the '-active' flag is passed
	case *active:
		output := ""
		// add a prefix to be displayed
		for idx, item := range *l {
			// display only active tasks
			if !item.Done {
				prefix := "[ ] "
				output += fmt.Sprintf("%s%d: %s\n", prefix, idx+1, item.Task)
			}
		}
		fmt.Fprint(os.Stdout, output)

	// INFO: check case where the '-verbose' flag is passed
	case *verbose:
		status := "Active"
		output := ""
		prefix := "[ ] "
		for idx, item := range *l {
			timeString := item.CreatedAt.Format(time.UnixDate)
			if item.Done {
				status = "Done"
				prefix = "[x] "
			}
			output += fmt.Sprintf("%s%d: %s | Created: %s | Status: %s\n", prefix, idx+1, item.Task, timeString, status)
		}
		fmt.Fprint(os.Stdout, output)

	// check for the case where the '-list' flag is passed
	case *list:
		// Print(l) uses the default String() for the type, which in our case uses a for-range
		fmt.Print(l)

		// check for the case where the '-complete' flag is passed with positive value
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

		// check for the case where the '-add' flag is passed
	case *add:
		// if true, well call the getTask() method with the os.Stdin (which implements io.Reader).
		// for the variadic paramenter, pass the flag.Args() wich collects all non flag arguments passed to the command line.
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		// call Add() with the string getTasks returns
		l.Add(t)

		// save the updated list on disk.
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		//INFO: check case where the '-del' flag is passed with a positive value
	case *del > 0:
		// calls the Delete() method with the pos of the -del value
		if err := l.Delete(*del); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		//INFO: save the updated on disk
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
