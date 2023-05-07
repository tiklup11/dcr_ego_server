package utils

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InsertCSVToMySQL2(csvFilePath string) error {
	// Open the CSV file for reading
	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "<user>:<password>@tcp(<host>:<port>)/<database>")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new CSV reader
	csvReader := csv.NewReader(csvFile)

	// Read the CSV header
	header, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("CSV file is empty")
		}
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	// Prepare the SQL statement for inserting data
	sqlStmt := "INSERT INTO <table> (" + header[0]
	for _, colName := range header[1:] {
		sqlStmt += ", " + colName
	}
	sqlStmt += ") VALUES (?,"
	for i := 1; i < len(header); i++ {
		sqlStmt += "?,"
	}
	sqlStmt = sqlStmt[:len(sqlStmt)-1] + ")"

	// Prepare the SQL statement for executing
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Read each row of the CSV file and insert it into the database
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read CSV row: %v", err)
		}

		// Convert the row values to interface{} for passing to the SQL statement
		var rowVals []interface{}
		for _, val := range row {
			rowVals = append(rowVals, val)
		}

		// Execute the SQL statement with the row values
		_, err = stmt.Exec(rowVals...)
		if err != nil {
			return fmt.Errorf("failed to insert row into database: %v", err)
		}
	}

	return nil
}

