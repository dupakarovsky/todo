package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

//====================
// TESTING
//====================
/*
   Testing will need to:
   > use the go build tool to compile the program into a binary
   > execute the binary with different arguments and test the correct behaviour

    use the TestMain function to help control extra tasks required to setup and tear down the resources necessary for testing.
*/

// variables will hold the binary name and the file name, which is required for the Save() method to be called.
// will be used by all tests in this file
var (
	binName  = "todo"
	fileName = "todo.json"
)

// TestMain will test wether we're able to compile the todo code into an executable binary.
func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	//INFO: use the TODO_FILENAME environmental variable
	if os.Getenv("TODO_FILENAME") != "" {
		fileName = os.Getenv("TODO_FILENAME")
	}

	// access the GOOS variable during runtime to check if we're builing for windows, and, if so, add the .exe extension.
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// construct the command to build the binary with with "go", with args: 'build', '-o'
	// and the binary name we've give (eg.: go build -o todo.exe)
	build := exec.Command("go", "build", "-o", binName)

	// try running the build command
	if err := build.Run(); err != nil {
		// On error, print the error an exit
		fmt.Fprintf(os.Stderr, "error building binary: %s", err.Error())
		os.Exit(1)
	}

	// build successful, run any test function in this file
	fmt.Println("Running tests...")

	// store the exit code the binary has returned
	code := m.Run()

	// clean up the binary file and .todo.json file
	fmt.Println("Cleaning up..")
	os.Remove(binName)
	os.Remove(fileName)

	// exit with the returned code
	os.Exit(code)

}

// TestTodoCLI will use the subtest feature to execute tests that depend on each other with
// the t.run() method
func TestTodoCLI(t *testing.T) {
	// defines a new task name
	task := "test task number 1"

	// get the current working directory. Abort if it fails
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// store the path to the binary compiled in TestMain()
	cmdPath := filepath.Join(dir, binName)

	// Create a subtest (adds a new task)
	t.Run("AddNewTask", func(t *testing.T) {
		// build a command to run: the path to our compiled binary CLI and the arguments we're passing.
		// here we're splitting the 'task' string into multiple strings and passing them as arguments to the CLI binary
		// (eg.: /mnt/dev/WebDev/go/powerclag/020_user_interaction/todo/cmd/todo test task number #1)

		// update the command to run with the '-add' flag instead. the task string will be passed first non flag argument
		// to the command
		cmd := exec.Command(cmdPath, "-add", task)

		// run the command
		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// AddNewTaskStdIN subtest will test for the case where the '-add' flag is passed, but no arguments. This will trigger
	// the getTasks() function to read from the StdIn. the test case will use the cmd.StdinPipe() to read from the commands standard
	// input and the io.WriteString will act as the user typing, by passing the task2 string to the cmdStdIn pipe
	task2 := "test task number 2"
	t.Run("AddNewTaskStdIN", func(t *testing.T) {
		//build the command with only the '-add' flag, without any non-flag arguments.
		cmd := exec.Command(cmdPath, "-add")

		// starts a pipe connected to the command stdin when the command starts starts
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		// io.WriteString will write the task2 string to the cmdStdIn pipe
		io.WriteString(cmdStdIn, task2)

		// close the stream after the string is written.
		cmdStdIn.Close()

		// run the command
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Create a subtest (List tasks).
	t.Run("ListTasks", func(t *testing.T) {
		// build the command to run: path to the compiled binary without any arguments. all strings after 'todo' will be concatenated
		// into
		cmd := exec.Command(cmdPath, "-list")

		// run the command and return the combined std output and stderr
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// perform an assertion on what to expect.
		expected := fmt.Sprintf("[ ] 1: %s\n[ ] 2: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("expected %q; got %q instead\n", expected, string(out))
		}
	})

	//INFO: Create a subtest (CompleteTask);
	t.Run("CompleteTask", func(t *testing.T) {

		// build a command to run: path to the compiled binary with the -complete flag and 1, to complete the 1st task
		cmd := exec.Command(cmdPath, "-complete", "1")

		// run the command
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	//INFO: Create a subtest (DeleteTask);
	t.Run("DeleteTask", func(t *testing.T) {

		// build a command to run: path to the compiled binary with the -del flag and 1 to delete the 1st task
		cmd := exec.Command(cmdPath, "-del", "1")

		// run the command
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

}

// ===============================
// CLEAR
// ===============================

// > Clear the TODO_FILENAME environmental variable and remove the todo.json/new-todo.json file
// $ unset TODO_FILENAME
// $ rm new-todo.json

// > Rebuild the binary
// $ go build -o=./bin/todo ./cmd

// > Run the app with some none flag arguments
// $ ./todo -add Incluing item from Args
// $ ./todo -list  // [ ] 1: Incluing item from Args

// > Pipe a string from stdin into the command
// $ echo 'This item comes from STDIN' | ./todo -add
