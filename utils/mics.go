package utils

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func ChangeDirTo(path string) {
	// Change the working directory to the specified path
	err := os.Chdir(path)
	fmt.Print("changing directory to tranformation_code :", path)
	fmt.Print("\n")
	if err != nil {
		fmt.Print(err)
	}
}

func DeleteFolder(tempFolderPath string) {
	err := os.RemoveAll(tempFolderPath)
	if err != nil {
		fmt.Errorf("not able to remove %s", tempFolderPath)
	}
}

func InsertCsvToTable(db *sql.DB, filePath string, tableName string) error {
	// Open the CSV file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1 // allow variable number of fields
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	// Build the SQL create table statement
	var columnDefs []string
	for _, header := range headers {
		columnDefs = append(columnDefs, fmt.Sprintf("%s TEXT", header))
	}
	createTableStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columnDefs, ", "))

	// Create the table if it doesn't already exist
	_, err = db.Exec(createTableStmt)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	// Build the SQL insert statement
	numColumns := len(headers)
	placeholders := strings.Repeat("?, ", numColumns)
	placeholders = strings.TrimSuffix(placeholders, ", ")
	insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(headers, ", "), placeholders)

	// Insert each record into the MySQL table
	var values []interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %v", err)
		}

		for _, value := range record {
			values = append(values, value)
		}

		_, err = db.Exec(insertStmt, values...)
		if err != nil {
			return fmt.Errorf("failed to insert record into table: %v", err)
		}

		values = values[:0] // reset values
	}

	return nil
}
