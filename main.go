package main

import (
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"strconv"
	"strings"
)

func main() {

	config := DBConfig{
		Server:   "sql-alp-archqst-001.database.windows.net", // "139.53.211.202",
		Port:     0,                                          // 1433
		Database: "QSTARCHIVES",                              // "QST_Dev",
		AppName:  "autoTransfer",
	}

	conn, err := NewDBConn(config)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	ok, errPing := conn.Ping()
	if errPing != nil {
		fmt.Println(errPing)
	} else {
		fmt.Println(ok)
	}

	// query := "SELECT * FROM Archive.MeasurementAggregation WHERE Timestamp BETWEEN '2023-12-01' AND '2023-12-02'"
	query := "SELECT TOP 1 * FROM dbo.MeasurementAggregations ORDER BY Id desc"
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Query error:", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Error getting columns:", err)
	}

	// Create a slice of interface{} to hold values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Print column headers
	fmt.Println(strings.Join(columns, " | "))

	// Iterate over rows
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			fmt.Println("Row scan error:", err)
		}

		var rowData []string

		for _, val := range values {
			switch v := val.(type) {
			case nil:
				rowData = append(rowData, "NULL")
			case []byte:
				// Convert byte slice to string
				rowData = append(rowData, string(v))
			case float64:
				// Format float with precision
				rowData = append(rowData, strconv.FormatFloat(v, 'f', 6, 64))
			default:
				// Convert everything else to string
				rowData = append(rowData, fmt.Sprintf("%v", v))
			}
		}

		// fmt.Println(strings.Join(rowData, " | "))
		fmt.Println(rowData)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Rows iteration error:", err)
	}

}
