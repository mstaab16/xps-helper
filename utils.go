package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// //go:embed transformed_data.csv
// var dataCSV embed.FS

type DataRow struct {
	Number  int
	Element string
	Orbital string
	Energy  float64
}

func loadCSV() ([]DataRow, error) {
	file, err := dataCSV.ReadFile("transformed_data.csv")
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bytes.NewReader(file))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []DataRow
	for _, record := range records[1:] { // Skip header
		number, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		energy, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		data = append(data, DataRow{
			Number:  number,
			Element: record[1],
			Orbital: record[2],
			Energy:  energy,
		})
	}
	return data, nil
}

func filterByEnergy(data []DataRow, minEnergy, maxEnergy float64) []DataRow {
	var result []DataRow
	for _, row := range data {
		if row.Energy >= minEnergy && row.Energy <= maxEnergy {
			result = append(result, row)
		}
	}
	return result
}

func filterByElement(data []DataRow, name string) []DataRow {
	var result []DataRow
	for _, row := range data {
		if strings.ToLower(row.Element) == strings.ToLower(name) {
			result = append(result, row)
		}
	}
	return result
}

// isAlpha checks if a string contains only letters (A-Z, a-z)
func isAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0 // Ensures the string is not empty
}
func parseFloatQuery(input string) (float64, float64, error) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid input format: %s", input)
	}

	// Parse the first float
	x, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing first number: %w", err)
	}

	// Parse the second float
	y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing second number: %w", err)
	}

	return x, y, nil
}

// func main() {
// 	data, err := loadCSV()
// 	if err != nil {
// 		log.Fatalf("Failed to load CSV: %v", err)
// 	}

// 	filtered := filterByEnergy(data, -5.0, 5.0) // Example bounds
// 	for _, row := range filtered {
// 		fmt.Printf("%d, %s, %s, %.2f\n", row.Number, row.Element, row.Orbital, row.Energy)
// 	}

// 	filtered_ele := filterByElement(data, "K") // Example bounds
// 	for _, row := range filtered_ele {
// 		fmt.Printf("%d, %s, %s, %.2f\n", row.Number, row.Element, row.Orbital, row.Energy)
// 	}
// }
