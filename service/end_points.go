package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tiklup11/go-hello-world/constants"
	"github.com/tiklup11/go-hello-world/database"
)

func isTranformationAlreadyDone(refid string) bool {
	fmt.Print("checking if tranformation is already executed")
	fmt.Print("\n")

	// return err != nil
	value, err := database.Read(refid, constants.DbPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("tranformation is already executed")
	fmt.Println("value for refid ", refid, " = ", value)
	return true
}

// this is to handle [/Run POST] requests
func RunHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("POST request recieved at /run \n")
	fmt.Print("processing request.. \n")
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Print(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(32 << 20) // Max file size 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Print(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the value of the myParam field
	refId := r.FormValue("refid")
	fmt.Print()

	//checking if the solution already exists
	isAlreadyDone := isTranformationAlreadyDone(refId)

	if isAlreadyDone {
		fmt.Fprintf(w, "tranformation already executed, to get results, request /download with refid")
		return
	}

	// Get the ZIP file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary directory to extract the files from the ZIP
	tmpDir, err := os.MkdirTemp("", "zip-extract")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Print(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	// Extract the files from the ZIP to the temporary directory
	zipReader, err := zip.NewReader(file, handler.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, f := range zipReader.File {
		fpath := filepath.Join(tmpDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(tmpDir)+string(os.PathSeparator)) {
			http.Error(w, "Invalid ZIP file", http.StatusBadRequest)
			return
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		inFile, err := f.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer inFile.Close()
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Save the extracted files to the server directory
	destDir := "../go_tranformations"
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = os.Rename(tmpDir, filepath.Join(destDir, handler.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// running script
	_ = GenerateScriptAndRun(refId)

	// fmt.Print("logging output..\n")
	// fmt.Print(output)

	// fmt.Fprintf(w, output)
}

// not in use, endpoint : [/ GET]
func HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Read the entire file into memory

	refid := r.URL.Query().Get("refid")
	if refid == "" {
		http.Error(w, "refid parameter is missing", http.StatusBadRequest)
		return
	}

	value, err := database.Read(refid, constants.DbPathFromOutside)

	if err != nil {
		value, err = database.Read(refid, constants.DbPath)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}

	data := map[string]string{
		refid: value,
	}

	// Marshal the map into a JSON byte slice
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Do something with the refid parameter
	fmt.Fprintf(w, string(jsonData))
}

// not in use, endpoint : [/ GET]
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, HTTPS World!")
}
