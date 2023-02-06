package main

// add BODY -> csv file
// open csv file

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var CSVs [][][]string
	ch := make(chan [][]string)
	startTime := time.Now()
	outFile := "data/out.csv"
	URLS := os.Args[1:]

	fmt.Printf("\n+%90s+\n", strings.Repeat("-", 90))
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s", "FILE NAME", "PROCESS", "TOTAL TIME", "PROCESS TIME", "DETAILS")
	fmt.Printf("\n+%90s+\n", strings.Repeat("-", 90))
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | Threads utilized: %d\n", "", "STARTING", time.Since(startTime).String(), "", len(URLS))

	// Read file, translate to CSV, clean and sort
	for i := 0; i < len(URLS); i++ {
		wg.Add(1)
		go runCSVJobs(URLS[i], ch, startTime, &wg)
	}

	for i := 0; i < len(URLS); i++ {
		CSVs = append(CSVs, <-ch)
	}

	// WaitGroup jobs for concurrency
	wg.Wait()
	close(ch)

	// Merge CSVs
	var sortedCSV [][]string
	for i := 0; i < len(CSVs); i++ {
		curCSV := CSVs[i]
		if len(curCSV) > 1 {
			sortedCSV = append(sortedCSV, curCSV...)
		}
	}

	sortedCSV = sortCSV(sortedCSV, startTime, outFile)

	// Finds mean and median
	printStats(sortedCSV, startTime, &wg)
	wg.Wait()
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | \n", outFile, "FINISHED", time.Since(startTime).String(), "")

	// Writes data to CSV
	// toCSV := append([][]string{{"fname", "lname", "age"}}, sortedCSV...)
	// writeToCSV(toCSV, outFile, startTime)
}

// Runs the set of functions to clean and sort a CSV
func runCSVJobs(url string, ch chan [][]string, startTime time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	var sortedCSV [][]string
	isValid := validateURL(url, startTime)
	var filePath string
	if isValid == 1 {
		filePath = strings.Trim(url, "file://")
		csv, check := readCSVFile(filePath, startTime)
		if check == true {
			cleanedCSV := cleanCSV(csv, startTime, filePath)
			sortedCSV = sortCSV(cleanedCSV, startTime, filePath)
		} else {
			sortedCSV = [][]string{{"fname", "lname", "age"}}
		}
	} else {
		sortedCSV = [][]string{{"fname", "lname", "age"}}
	}
	ch <- sortedCSV
}

// input CSV is valid
// Reads CSV file from file location to output matrix [][]string and bool stating if the
func readCSVFile(filePath string, startTime time.Time) ([][]string, bool) {
	func_time := time.Now()
	var records [][]string
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), "File cannot be read")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err = csvReader.ReadAll()
	if err != nil {
		fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), err.Error())
	} else if len(records) > 0 {
		if reflect.DeepEqual(records[0], []string{"fname", " lname", " age"}) {
			fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "READING", time.Since(startTime).String(), time.Since(func_time).String(), "")
			return records, true
		} else {
			fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), "Has the incorrect CSV headers: ["+strings.Join(records[0], ", ")+"]")
		}
	} else {
		fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), "File is empty")
	}
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "READING", time.Since(startTime).String(), time.Since(func_time).String(), "")
	return [][]string{{"fname", "lname", "age"}}, false
}

// Checks if URL Path is Valid by file://
// 0: "Not valid URL"
// 1: "file://"
// 2: "https://"
func validateURL(filePath string, startTime time.Time) int {
	func_time := time.Now()
	if strings.Contains(filePath, "file://") {
		return 1
	} else if strings.Contains(filePath, "http://") {
		return 2
	} else {
		// do stuff for non local files (https://, ftp://, etc.)
		// ex. download data from online database using the package net/http into ./data/[FILE]
		fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), "File is not properly formatted under file ext. 'file://'")
		return 0
	}
}

// Returns a CSV with cells with missing items removed
func cleanCSV(csv [][]string, startTime time.Time, filePath string) [][]string {
	func_time := time.Now()
	var cleanedCSV [][]string
	for i := 1; i < len(csv); i++ {
		fname := strings.Trim(csv[i][0], " ")
		lname := strings.Trim(csv[i][1], " ")
		age := strings.Trim(csv[i][2], " ")
		if fname != "" && lname != "" && age != "" {
			cleanedCSV = append(cleanedCSV, csv[i])
		}
	}
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %d cell[s] removed \n", filePath, "CLEANING", time.Since(startTime).String(), time.Since(func_time).String(), len(csv)-len(cleanedCSV))
	return cleanedCSV
}

// Sort element (age) of a CSV
func sortCSV(csv [][]string, startTime time.Time, filePath string) [][]string {
	func_time := time.Now()
	sort.SliceStable(csv, func(i, j int) bool {
		age1 := strings.Trim(csv[i][2], " ")
		age2 := strings.Trim(csv[j][2], " ")
		firstNum, err := strconv.Atoi(age1)
		if err != nil {
			fmt.Printf("| %-40s | %-15s | %-12s | %-12s | \n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String())
		}
		secondNum, err := strconv.Atoi(age2)
		if err != nil {
			fmt.Printf("| %-40s | %-15s | %-12s | %-12s | \n", filePath, "ERROR", time.Since(startTime).String(), time.Since(func_time).String())
		}
		return firstNum < secondNum
	})
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | \n", filePath, "SORTING", time.Since(startTime).String(), time.Since(func_time).String())
	return csv
}

// Concurrently runs functions `printMedian` and `printMean`
func printStats(sortedCSV [][]string, startTime time.Time, wg *sync.WaitGroup) {
	wg.Add(2)
	go printMedian(sortedCSV, startTime, wg)
	go printMean(sortedCSV, startTime, wg)
}

// Finds median of CSV
func printMedian(csv [][]string, startTime time.Time, wg *sync.WaitGroup) {
	func_time := time.Now()
	defer wg.Done()
	csv = csv[1:]
	var median float64
	l := len(csv)
	medianPerson := csv[l/2]
	if l == 0 {
		median = 0
	} else if l%2 == 0 {
		firstNum, _ := strconv.ParseFloat(strings.Trim(csv[l/2-1][2], " "), 64)
		secondNum, _ := strconv.ParseFloat(strings.Trim(csv[l/2][2], " "), 64)
		median = (firstNum + secondNum) / 2
	} else {
		median, _ = strconv.ParseFloat(strings.Trim(csv[l/2][2], " "), 64)
	}
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %.2f\n", "RESULT", "MEDIAN", time.Since(startTime).String(), time.Since(func_time).String(), median)
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %s\n", "RESULT", "MEDIAN PERSON", time.Since(startTime).String(), time.Since(func_time).String(), medianPerson)
}

// Prints mean
func printMean(csv [][]string, startTime time.Time, wg *sync.WaitGroup) {
	func_time := time.Now()
	defer wg.Done()
	result := 0
	var newNum int
	var mean float64
	var err error
	for i := 0; i < len(csv); i++ {
		newNum, err = strconv.Atoi(strings.Trim(csv[i][2], " "))
		if err != nil {
			fmt.Printf("| %-35s | %-15s | %-12s | %-12s | %s\n", "RESULT", "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), err)
		}
		result += newNum
	}
	mean = float64(result) / float64(len(csv))
	fmt.Printf("| %-40s | %-15s | %-12s | %-12s | %.2f\n", "RESULT", "MEAN", time.Since(startTime).String(), time.Since(func_time).String(), mean)
}

// Writes [][]string To CSV
// Extra function to show final product
func writeToCSV(data [][]string, fileName string, startTime time.Time) {
	func_time := time.Now()
	csvFile, err := os.Create(fileName)

	if err != nil {
		log.Fatalf("| %-35s | %-15s | %-12s | %-12s | %s\n", fileName, "ERROR", time.Since(startTime).String(), time.Since(func_time).String(), "Error creating file"+fileName)
	}
	csvwriter := csv.NewWriter(csvFile)

	for _, row := range data {
		row[0] = strings.Trim(row[0], " ")
		row[1] = strings.Trim(row[1], " ")
		row[2] = strings.Trim(row[2], " ")
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvFile.Close()
	fmt.Printf("| %-35s | %-15s | %-12s | %-12s | \n", fileName, "WRITING", time.Since(startTime).String(), time.Since(func_time).String())
}
