package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tiklup11/go-hello-world/database"
	"github.com/tiklup11/go-hello-world/utils"
)

func GenerateScriptAndRun(refId string) string {

	go_tranformation_path := "../go_tranformations/test.zip"
	utils.ChangeDirTo(go_tranformation_path)

	// Define the script content as a string
	script := generateScript()

	// Write the script to a file
	err := os.WriteFile("myscript.sh", []byte(script), 0755)
	if err != nil {
		fmt.Println("Error writing script:", err)
		return "error writing script"
	}

	// Make the script executable
	cmd := exec.Command("chmod", "+x", "myscript.sh")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error making script executable:", err)
		return "error making script executable"
	}

	hasScriptFinished := make(chan bool)
	out := make(chan string)

	// Start a goroutine to execute RunScript in the background
	go func() {
		// Run the script and capture the output
		output := RunScript("./myscript.sh")
		// RunScript("./hello.sh")
		// Signal that the function has completed executing
		hasScriptFinished <- true
		out <- output
	}()

	// Wait for the RunScript function to complete before executing f2()
	<-hasScriptFinished
	output := <-out

	// Print the output
	fmt.Println(string(output))

	formatedOutput := formatOutput(string(output))
	fmt.Println("saving to db..")
	//saving to db
	saveToDB(refId, formatedOutput)

	utils.DeleteFolder(go_tranformation_path)

	return string(output)
}

func formatOutput(output string) string {
	lines := strings.Split(output, "\n")
	lastLine := lines[len(lines)-2]

	return lastLine
}

func saveToDB(key string, value string) {
	if err := database.InsertFunc(key, value); err != nil {
		fmt.Println(err)
		return
	}
}

func generateScript() string {
	script := `
#!/bin/bash

echo "Starting script to run go-tranformation.."

echo "running : [go mod tidy] to load dependencies of tranformation"
go mod tidy

echo "running : [ego-go build main.go]"
ego-go build main.go

ego sign main

echo "running : [modifying enclave.json config file] to include datafiles"
cp ../../ego_server/service/enclave.json ./ -f

ego sign main

echo "running : [OE_SIMULATION=1 ego run main] to set mode to SIMULATION"
OE_SIMULATION=1 ego run main
`
	return script
}

func RunScript(path string) string {
	cmd, err := exec.Command("/bin/sh", path).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string(cmd)
	return output
}

func TestScript() string {
	go_tranformation_path := "../go_tranformations/test.zip"
	utils.ChangeDirTo(go_tranformation_path)

	// Run the script and capture the output
	out := RunScript("./myscript.sh")
	return string(out)
}
