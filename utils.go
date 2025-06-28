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

// parseElementNames splits a string by commas and spaces, trims whitespace, and returns a slice of element names
func parseElementNames(input string) []string {
	// First split by commas
	commaParts := strings.Split(input, ",")
	var elements []string
	
	for _, part := range commaParts {
		// Then split each comma part by spaces
		spaceParts := strings.Fields(strings.TrimSpace(part))
		for _, element := range spaceParts {
			element = strings.TrimSpace(element)
			if element != "" {
				elements = append(elements, element)
			}
		}
	}
	
	return elements
}

func filterByElement(data []DataRow, name string) []DataRow {
	var result []DataRow
	elements := parseElementNames(name)
	
	// Convert element names to lowercase for case-insensitive comparison
	elementMap := make(map[string]bool)
	for _, element := range elements {
		elementMap[strings.ToLower(element)] = true
	}
	
	for _, row := range data {
		if elementMap[strings.ToLower(row.Element)] {
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

// isElementSearch checks if a search query is for elements (contains letters, spaces, commas)
// vs energy (should be a single number)
func isElementSearch(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	
	// If it contains any letters, it's an element search
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	
	// If it contains spaces or commas, it's likely an element search
	if strings.Contains(s, " ") || strings.Contains(s, ",") {
		return true
	}
	
	// Otherwise, it's probably an energy search (single number)
	return false
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
