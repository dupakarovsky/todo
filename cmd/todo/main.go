package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dupakarovsky/todo"
)

//=================================
// CAPTURING INPUT FROM STDIN
//=================================
// The abbility to add new tasks via stdin allow your users to pipe new tasks from other command-line tools.
// We'll use the 'bufio' package to read data form the stdin input stream.

// > Create a helper furnction called 'getTasks()' that accepts an parameter of a type that implements the io.Reader interface

// name of the json file that'll be created.
var todoFileName = "todo.json"

// getTask will acpet a first parameter that implements the io.Reader interface. Then a variadict string parameter to collect all
// others into a slice.
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
		flag.PrintDefaults()
	}

	//INFO: update the -task flag to -bool
	// Add a list of flags to be passed to the command line
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all ToDo items")
	complete := flag.Int("complete", 0, "Item to be completed")
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

	//  check for the case where the '-list' flag is passed
	case *list:
		// Pritnt(l) uses the default String() for the type, which in our case uses a for-range
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

		//INFO: check for the case where the '-add' flag is passed
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

		//TODO: > update the main_test.AddNewTaks test and also to create a new test

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
