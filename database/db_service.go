package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tiklup11/go-hello-world/constants"
)

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var databaseFilePath string = constants.DbPathFromOutside

// "./database/data.json"

func Read(key string, DbPath string) (string, error) {

	// Read the entire file into memory
	data, err := os.ReadFile(DbPath)
	if err != nil {
		// fmt.Print(err)
		return "", err
	}

	// Parse the JSON data into a slice of Pair structs
	var pairs []Pair
	if err := json.Unmarshal(data, &pairs); err != nil {
		return "", err
	}

	// Find the pair with the specified key
	for _, pair := range pairs {
		if pair.Key == key {
			return pair.Value, nil
		}
	}

	return "", fmt.Errorf("key not found: %s", key)
}

func pwd() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Current working directory:", cwd)
}

func InsertFunc(key string, value string) error {
	fmt.Print("inserting ", key, " : ", value)
	fmt.Print()

	// Read the entire file into memory
	data, err := os.ReadFile(databaseFilePath)
	if err != nil {
		return err
	}

	// Parse the JSON data into a slice of Pair structs
	var pairs []Pair
	if err := json.Unmarshal(data, &pairs); err != nil {
		return err
	}

	// Check if the key already exists
	for i, pair := range pairs {
		if pair.Key == key {
			pairs[i].Value = value
			newData, err := json.Marshal(pairs)
			if err != nil {
				return err
			}
			// Write the new data back to the file
			if err := os.WriteFile(databaseFilePath, newData, 0644); err != nil {
				return err
			}
			return nil
		}
	}

	// If the key doesn't exist, add a new pair to the slice
	pairs = append(pairs, Pair{Key: key, Value: value})

	// Encode the updated slice as JSON
	newData, err := json.Marshal(pairs)
	if err != nil {
		return err
	}

	// Write the new data back to the file
	if err := os.WriteFile(databaseFilePath, newData, 0644); err != nil {
		return err
	}

	return nil
}

func Test() {
	// Example usage
	// if err := InsertFunc("helqqlo", "worqqqld"); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	value, err := Read("4ew54!@23e1311", databaseFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}
