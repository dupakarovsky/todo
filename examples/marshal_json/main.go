package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

func main() {

	// writeJSON will marshal
	writeJSON := func(slice []string) error {

		// marshal the slice into json
		js, err := json.Marshal(slice)
		if err != nil {
			return fmt.Errorf("error marshaling json: %v", err)
		}

		// write to the os.
		err = os.WriteFile("mashalled.json", js, fs.FileMode(os.O_WRONLY))
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}

		fmt.Println("File written successfuly")
		return nil
	}

	// call the function
	err := writeJSON([]string{"Add New Task Master"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
