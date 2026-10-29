package main

// gpuinfo binary from "github.com/andersrennermalm/gpuinfo"
import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

func getGpu() {
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

func getMemPressure() (memPressure float64, err error) {

	memPressure = 0.0

	cmd := exec.Command("memory_pressure", "-Q")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return
	}
	logrus.Infof("Received output %s", out.String())

	// Regex to extract the "Device Utilization %" integer value
	re := regexp.MustCompile(`free percentage: (\d+)%`)
	matches := re.FindStringSubmatch(out.String())
	logrus.Infof("matches output %v", matches)
	if len(matches) < 2 {
		err = fmt.Errorf("could not parse memory pressure; confirm this is an Apple Silicon Mac")
		return
	}

	// Convert parsed string back to an float
	memPressure, err = strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}
	memPressure = 100 - memPressure
	logrus.Infof("memPressure Usage is  : %2.2f ", memPressure)

	return

}

func getCpu() (cpuUsage float64, err error) {

	cpuUsage = 0.0
	cmd := exec.Command("top", "-l", "1", "-n", "0", "-s", "1")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return
	}
	// logrus.Infof("Received output %s", out.String())
	// Regex to extract the "Device Utilization %" integer value
	re := regexp.MustCompile(`CPU usage:.*, (\d+.\d+)% idle`)
	matches := re.FindStringSubmatch(out.String())

	if len(matches) < 2 {
		err = fmt.Errorf("could not parse CPU utilization; confirm this is an Apple Silicon Mac")
		return
	}

	// Convert parsed string back to an float
	cpuUsage, err = strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}
	cpuUsage = 100 - cpuUsage
	logrus.Infof("CPU Usage is  : %2.2f ", cpuUsage)

	return
}

func main() {
	cpu, err := getCpu()
	mem, err := getMemPressure()
	if err != nil {
		logrus.Errorf("Cannot get metrics %v", err)
	}

	logrus.Printf("CPU %2.2f, MEM %2.2f", cpu, mem)
}
