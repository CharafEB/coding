package dawonlod

import (
    "encoding/csv"
    "os"
)

// WriteToCSV takes data as [][]string and writes it to a CSV file.
func WriteToCSV(filename string, data [][]string) error {
    // Create or open the file
    file, err := os.Create(filename)
	if err != nil {
        return err
    }
    defer file.Close()

    // Create a new CSV writer
    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write each row to the CSV file
    for _, row := range data {
        if err := writer.Write(row); err != nil {
            return err
        }
    }

    return nil
}