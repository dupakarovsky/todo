package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	//
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
		cmd := exec.Command(cmdPath, strings.Split(task, " ")...)

		// run the command
		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	// Create a subtest (List tasks).
	t.Run("ListTasks", func(t *testing.T) {
		// build the command to run: path to the compiled binary without any arguments. all strings after 'todo' will be concatenated
		// into
		cmd := exec.Command(cmdPath)

		// run the command and return the combined std output and stderr
		// ERROR: This test will fail with all the Println in main.
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// perform an assertion on what to expect.
		expected := task + "\n"
		if expected != string(out) {
			t.Errorf("expected %q; got %q instead\n", expected, string(out))
		}
	})

}
