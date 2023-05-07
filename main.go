package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/tiklup11/go-hello-world/service"
)

func main() {
	service.RunEgoServer()
}

// utils.DeleteFolder("./../go_transformations/test.zip/")

// func RunHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Print("POST request recieved at /run \n")
// 	fmt.Print("processing request.. \n")
// 	// Check if the request method is POST
// 	if r.Method != "POST" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		fmt.Print(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse the multipart form
// 	err := r.ParseMultipartForm(32 << 20) // Max file size 32MB
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		fmt.Print(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Get the ZIP file from the form
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		fmt.Print(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a temporary directory to extract the files from the ZIP
// 	tmpDir, err := os.MkdirTemp("", "zip-extract")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		fmt.Print(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	// Extract the files from the ZIP to the temporary directory
// 	zipReader, err := zip.NewReader(file, handler.Size)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		fmt.Print(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	for _, f := range zipReader.File {
// 		fpath := filepath.Join(tmpDir, f.Name)
// 		if !strings.HasPrefix(fpath, filepath.Clean(tmpDir)+string(os.PathSeparator)) {
// 			http.Error(w, "Invalid ZIP file", http.StatusBadRequest)
// 			return
// 		}
// 		if f.FileInfo().IsDir() {
// 			os.MkdirAll(fpath, os.ModePerm)
// 			continue
// 		}
// 		err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		defer outFile.Close()
// 		inFile, err := f.Open()
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		defer inFile.Close()
// 		_, err = io.Copy(outFile, inFile)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// Save the extracted files to the server directory
// 	destDir := "../go_tranformations"
// 	err = os.MkdirAll(destDir, os.ModePerm)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = os.Rename(tmpDir, filepath.Join(destDir, handler.Filename))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Return a success message
// 	output := generateScriptAndRun()

// 	fmt.Print(output)

// 	fmt.Fprintf(w, output)

// }

// func handleIndex(w http.ResponseWriter, r *http.Request) {
// 	html := `<html>
// 				<body>
// 					<form action="/run" method="post" enctype="multipart/form-data">
// 						<input type="file" name="file"><br><br>
// 						<input type="submit" value="Run">
// 					</form>
// 				</body>
// 			</html>`
// 	fmt.Fprint(w, html)
// }

// func changeDirTo(path string) {
// 	// Change the working directory to the specified path
// 	err := os.Chdir(path)
// 	fmt.Print("changing directory to tranformation_code :", path)
// 	fmt.Print("\n")
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// }

// func testScript() string {
// 	go_tranformation_path := "../go_tranformations/test.zip"
// 	changeDirTo(go_tranformation_path)

// 	// Run the script and capture the output
// 	out := runScript("./myscript.sh")
// 	return string(out)
// }

// func runScript(path string) string {
// 	cmd, err := exec.Command("/bin/sh", path).Output()
// 	if err != nil {
// 		fmt.Printf("error %s", err)
// 	}
// 	output := string(cmd)
// 	return output
// }

// func generateScriptAndRun() string {

// 	go_tranformation_path := "../go_tranformations/test.zip"
// 	changeDirTo(go_tranformation_path)

// 	// Define the script content as a string
// 	script := generateScript()

// 	// Write the script to a file
// 	err := os.WriteFile("myscript.sh", []byte(script), 0755)
// 	if err != nil {
// 		fmt.Println("Error writing script:", err)
// 		return "error writing script"
// 	}

// 	// Make the script executable
// 	cmd := exec.Command("chmod", "+x", "myscript.sh")
// 	err = cmd.Run()
// 	if err != nil {
// 		fmt.Println("Error making script executable:", err)
// 		return "error making script executable"
// 	}

// 	// Run the script and capture the output
// 	out := runScript("./myscript.sh")

// 	// Print the output
// 	fmt.Println(string(out))

// 	deleteTempFolder("../go_tranformations/test.zip")

// 	return string(out)
// }

// func generateScript() string {
// 	script := `
// #!/bin/bash

// echo "Starting script to run go-tranformation.."

// echo >> "running : [go mod tidy] to load dependencies of tranformation"
// go mod tidy

// echo >> "running : [ego-go build main.go]
// ego-go build main.go

// ego sign main

// echo >> "running : [modifying enclave.json config file] to include datafiles"
// cp ../../ego_server/enclave.json ./ -f
// ego sign main

// echo >> "running : [OE_SIMULATION=1 ego run main] to set mode to SIMULATION"
// OE_SIMULATION=1 ego run main
// `
// 	return script
// }

// func deleteTempFolder(tempFolderPath string) {
// 	// testFolderPath := "../go_tranformations/test.zip"
// 	err := os.RemoveAll(tempFolderPath)
// 	if err != nil {
// 		fmt.Errorf("not able to remove %s", tempFolderPath)
// 	}
// }

// // ------------------- CODE NOT IN USE-----------------------

// func runCommands(commands []*exec.Cmd, dir string) (string, error) {
// 	// Change the working directory to the specified path
// 	err := os.Chdir(dir)
// 	fmt.Print("changing directory to tranformation_code :", dir)
// 	fmt.Print("\n")
// 	if err != nil {
// 		return "", err
// 	}
// 	setEnvVariable()

// 	// Execute the commands in the specified directory
// 	var output string
// 	for _, cmd := range commands {
// 		fmt.Print("running : ")
// 		fmt.Print(cmd.String())
// 		fmt.Print("\n\n")
// 		out, err := cmd.Output()

// 		if err != nil {
// 			return "", err
// 		}
// 		// if i == len(commands)-1 {
// 		output += string(out)
// 		// }
// 	}
// 	return output, nil
// }

// func setEnvVariable() {
// 	// Get the current environment variables
// 	env := os.Environ()

// 	// Add a new environment variable to the slice
// 	env = append(env, "OE_SIMULATION=1")

// 	fmt.Println("Simulation mode on, OE_SIMULATION set to true")
// }

// func runTranformationGoCode() string {

// 	fmt.Print("running tranformation..\n")

// 	dir := "../go_tranformations/test.zip"

// 	cmd_tidy := exec.Command("go", "mod", "tidy")
// 	cmd_build := exec.Command("ego-go", "build", "main.go")
// 	cmd_sign := exec.Command("ego", "sign", "main")
// 	// cmd_source_bashrc := exec.Command("source", "~/.bashrc")
// 	cmd_generate_enclave := exec.Command("cp", "../../ego_server/enclave.json", "./", "-f")

// 	// cmd_set_to_simulation := exec.Command("/bin/sh", "-c", "export OE_SIMULATION=1")
// 	cmd_run := exec.Command("ego", "run", "main")

// 	commands := []*exec.Cmd{cmd_tidy, cmd_build, cmd_sign, cmd_generate_enclave, cmd_run}

// 	output, err := runCommands(commands, dir)
// 	if err != nil {
// 		fmt.Println("Error running commands:", err)
// 		return ""
// 	}
// 	fmt.Println()
// 	fmt.Println("Output:", output)
// 	// deleteTempFolder("../go_tranformations/test.zip")
// 	return output
// }
