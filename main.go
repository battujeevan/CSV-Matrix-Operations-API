package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	fmt.Println("Csv file operations-Echo, Invert, Sum, Transpose")
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/invert", transposeHandler)
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/mul", multiplyHandler)
	http.HandleFunc("/flatten", flattenHandler)
	log.Println("starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/************************************************************************/
// Echo Handler function:
func echoHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	matrix, err := parseCSVFile(file)
	if err != nil {
		http.Error(w, "Failed to parse file", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(matrixToString(matrix)))
}

// invert Handler function:
func transposeHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	matrix, err := parseCSVFile(file)
	if err != nil {
		http.Error(w, " Failed to parse CSV file", http.StatusBadRequest)
		return
	}
	transpose := transposeOfMatrix(matrix)
	w.Header().Set("Content-Type", "text/Plain")
	w.Write([]byte(matrixToString(transpose)))
}

// Sum Handler function:
func sumHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	matrix, err := parseCSVFile(file)
	if err != nil {
		http.Error(w, "Failed to Open file", http.StatusBadRequest)
	}
	totalSum := 0
	for _, row := range matrix {
		for _, value := range row {
			totalSum += value
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("sum:%d", totalSum)))
}

// Multiplication Handler function:
func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	matrix, err := parseCSVFile(file)
	if err != nil {
		http.Error(w, "Failed to Open file", http.StatusBadRequest)
	}
	totalMul := 1
	for _, row := range matrix {
		for _, value := range row {
			totalMul *= value
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("sum:%d", totalMul)))
}

// Flatten Handler function:
func flattenHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	matrix, err := parseCSVFile(file)
	if err != nil {
		http.Error(w, "Failed to parse CSV file", http.StatusBadRequest)
	}
	flattened := converToFlattenMatrix(matrix)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(flattnedMatrixToFlattnedString(flattened)))
}

/********************************************************************/
/* Helper Functions */
func parseCSVFile(r io.Reader) ([][]int, error) {
	readrer := csv.NewReader(r)
	var matrix [][]int

	for {
		record, err := readrer.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		var row []int
		for _, value := range record {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("error while parsing value :%s, %v", value, err)
			}
			row = append(row, intValue)
		}
		matrix = append(matrix, row)
	}
	return matrix, nil
}

func matrixToString(matrix [][]int) string {
	var builder strings.Builder
	for _, row := range matrix {
		strRow := make([]string, len(row))
		for i, val := range row {
			strRow[i] = strconv.Itoa(val)
		}
		builder.WriteString(strings.Join(strRow, ",") + "\n")
	}
	return builder.String()
}
func transposeOfMatrix(matrix [][]int) [][]int {
	if len(matrix) == 0 {
		return nil
	}
	transposed := make([][]int, len(matrix[0]))
	for i := range transposed {
		transposed[i] = make([]int, len(matrix))
	}
	for i := range matrix {
		for j := range matrix[0] {
			transposed[j][i] = transposed[i][j]
		}
	}
	return transposed
}
func converToFlattenMatrix(matrix [][]int) []int {
	var flattned []int
	for _, row := range matrix {
		flattned = append(flattned, row...)
	}
	return flattned
}
func flattnedMatrixToFlattnedString(flattened []int) string {
	row := make([]string, len(flattened))
	for i, val := range flattened {
		row[i] = strconv.Itoa(val)
	}
	return strings.Join(row, ",")
}

/********************************************************************/
