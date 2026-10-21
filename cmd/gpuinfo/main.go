package main

// gpuinfo binary from "github.com/andersrennermalm/gpuinfo"
import (
	"fmt"
	"os/exec"
)

func main() {
	// Create the command
	// Note: If "gpuinfo" requires arguments, add them here: exec.Command("gpuinfo", "-a")
	cmd := exec.Command("./gpuinfo -p")

	// Execute the command and capture the standard output
	output, err := cmd.Output()
	if err != nil {
		// This will catch if the command doesn't exist or fails to run
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	// The output is returned as a byte slice ([]byte)
	// You can convert it to a string to store it in a variable
	result := string(output)

	fmt.Println("Command Output:")
	fmt.Println(result)
}
