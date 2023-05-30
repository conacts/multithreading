# crowdstrike-test

## 0. Intro
To run the project, we use the bash script `run.sh`. The questions to run the go file are shown below.
```bash
$ sh run.sh 
Would you like to run the concurrent test or nonconcurrent test (c/n)? n
Would you like to pull files over http or locally (h/l)? l
```
- `n` runs nonconcurrent file
- `c` runs concurrent file
- `h` pulls files over http
- `l` pulls files from local file system

Here we host our files located in the `./data` directory using `python3 -m http.server`, a file hosting library in the python standard library. This enables us to pull our files from port 8000 on our own machine. This works well for prototyping, however it fails to serve files at the rate we request them causing the concurrent program (using http) to fail. This could be solved by either
1. Requesting the files until they are sent
2. Requesting them from a proper server online

This solution was built for testing purposes and doesn't fully represent the potential of the program. 

## 1. Dependencies:
All dependencies used in the program are in the golang standard library.

```go
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
```

## 2. Functions:

### 2.1 runCSVJobs:
A job running tasks on a file sequentially. This includes:
1. Validate input URL.
2. Read input file.
3. Clean the CSV.
4. Sort the CSV's rows by `age`.
```go
func runCSVJobs(url string, ch chan [][]string, startTime time.Time, wg *sync.WaitGroup) 
	...
}
```

###### Parameters:
- `url` :  `file://[file].csv`.
- `ch` : Channel to write to 
- `startTime` :  Used to measure function performance.
- `wg` : Tool used to help time concurrency and wait till jobs are finished

###### Return:
- `[][]string` : Returns a sorted CSV.

### 2.2 readCSVFile:
Reads file to transform into CSV ([stack overflow](https://stackoverflow.com/questions/24999079/reading-csv-file-in-go)).
```go
func readCsvFile(filePath string, startTime time.Time) ([][]string, bool) {
	...
}
```

###### Parameters:
- `filePath` : Command line parameter to translate file path to CSV.
- `startTime` :  Used to measure function performance.

###### Returns: 
- `[][]string` : CSV file.
- `bool` : Used to see if file was able to be read.

### 2.3 validateURL:
Checks if URL given from the command line is able to be parsed of form `file://`.
```go
func validateURL(filePath string, startTime time.Time) bool  {
	...
}
```

Note: The portion of the function dealing with https is left blank.

###### Parameters:
- `filePath` : Command line parameter to check if file is of the proper formatting.
- `startTime` :  Used to measure function performance.

###### Returns:
- `bool` : states if URL is valid or not.

### 2.4 cleanCSV:
Iterates through each element in CSV (in form of  `[][]string`) to remove elements with empty parameters. This problem could also have been accomplished through replacing data missing the `age` parameter with the mean of the data.
```go
func cleanCSV(csv [][]string, startTime time.Time, filePath string) [][]string {
	...
}
```

###### Parameters:
- `csv` : CSV file to be cleaned.
- `startTime` :  Used to measure function performance.
- `filePath` : Used to print function performance or errors.

###### Return:
- `[][]string` : A CSV with all the rows containing empty cells removed.

### 2.5 sortCSV:
Sorting CSV by age column.
```go 
func sortCSV(csv [][]string, startTime time.Time, filePath string) [][]string {
	...
}
```

###### Parameters
- `csv` : CSV file to be sorted.
- `startTime` :  Used to measure function performance.
- `filePath` : Used to print function performance or errors.

###### Returns:
- `[][]string` : Sorted CSV.

### 2.6 printStats:
Calls `printMedian(...)`  and `printMean(...)` functions concurrently.
```go
func printStats(sortedCSV [][]string, startTime time.Time, wg *sync.WaitGroup) {
	...
}
```

###### Parameters:
- `csv` : Desired CSV files passed into functions `printMedian` and `printMean`.
- `startTime` :  Used to measure function performance.
- `wg` : Tool used to help time concurrency and wait till jobs are finished

### 2.7 printMedian:
Prints median value and person containing median age datapoint to terminal.
```go
func printMedian(csv [][]string, startTime time.Time, wg *sync.WaitGroup) {
	...
}
```

###### Parameters:
- `csv` : Finds median and "median person" of a CSV file.
- `startTime` :  Used to measure function performance.
- `wg` : Tool used to help time concurrency and wait till jobs are finished

### 2.8 printMean:
Prints mean to terminal.
```go 
func printMean(csv [][]string, startTime time.Time, wg *sync.WaitGroup) {
	...
}
```

###### Parameters:
- `csv` : Finds mean of a CSV file.
- `startTime` :  Used to measure function performance.
- `wg` : Tool used to help time concurrency and wait till jobs are finished

### 2.9 downloadCSV
Downloads CSV from given url and stores it to filePath.
```go
func downloadCSV(url string, filePath string, startTime time.Time) bool {
	...
}
```

###### Parameters:
- `url` : URL to download CSV from
- `filePath` : File path to store CSV
- `startTime` :  Used to measure function performance.

###### Returns:
- `bool` : dictates whether CSV was downloaded or not

### 2.10 hashFilePath
A hashing function that takes in downloaded file name and current time to create a unique hash to prevent collisions.

```go
func hashFilePath(filePath string) string {
	...
}
```

###### Parameters:
- `filePath` : Used to hash a unique file name

###### Returns: 
- `string` : Unique file name


### 2.11 writeToFile:
A function to write bytes to a given file.
```go
func writeFile(data []byte, fileName string) bool {
	...
}
```


###### Parameters:
- `data` : Data to be written to file

###### Returns: 
- `bool` : Output whether write to file was successful or not

### 2.12 writeToCSV:
Writes cleaned, sorted and merged CSV to new file. Has fatal error if unable to write to file.
```go
func writeToCSV(data [][]string, fileName string) bool {
	...
}
```

###### Parameters:
- `data` : Data to write to CSV.
- `fileName` : Name of file we are writing to.

###### Return:
- `bool` : dictates whether saved to file or not.

## 3. Code Output:

##### Concurrent vs. Nonconcurrent
We receive a $44.37\%$ improvement using the concurrent file reading technique compared to the nonconcurrent version.

- Concurrent ~ $8.985542ms$
- Nonconcurrent ~ $12.97275ms$

$$\frac{12.97275ms - 8.985542ms}{8.985542ms} \cdot 100\% = 33.74\%$$

###### Concurrent Output: (Using local files)
![[Pasted image 20230206002719.png]]

###### Nonconcurrent Output: (Using local files)
![[Pasted image 20230206011617.png]]

## 5. Questions:

##### What assumptions did you make in your design? Why?
- Only one machine is present to read files. The addition of more machines could enable separate machines to work on files independently.
- I assumed that the machine had at least enough threads to support all the input files all at once.
- I assumed the machine had at least enough memory to handle all the files it was working on.
- I assumed that the CSVs being read are not always gonna be formatted properly and every input file needs to be checked.
- I assumed that if we were pulling over HTTP, we would be able to receive the file at each request. 


##### How would you change your program if it had to process many files where each file was over 10M records?

**With Multiple Machines:** 
The advantage of multiple machines is much more compute. This would allow us to be more sparing with the workload on each machine. We would also need to host the sorted CSVs for a master machine with lots of memory to merge into one final product. Multiple machines would also enable us to grab files based on distance from server we are communicating with for further optimization as well.
- I have created the `writeToCSV()` function to translate our golang CSVs back to regular CSV files to forward to master machine to work with.

**With a Single Machine:** 
With a single machine we will most likely be limited by the power of our machine. We will need ample RAM to store all the files as well as schedule how many files we work on at once. We should find the mean while sorting to prevent having to iterate through the data again. Another optimization could be storing our parsed data `age` as integers rather than strings.   This would prevent computing `strconv.Atoi(...)` and `strconv.ParseFloat(...)` more than once. The reason this wasn't implemented was because i would need to copy the entire array to a new array of form `[][]interface{}`

##### How would you change your program if it had to process data from more than 20K URLs?
- I would deploy more machines to have access to more compute.
- Have enough RAM on each machine to handle all possible file sizes
- If a requested file repeatedly times out, add it to the bottom of the list of files to work with to be used later
- Check requested file size so a machine can distribute workload to work on both larger and smaller files at the same time. This will enable the machine to more effectively use its compute.


##### How would you test your code for production use at scale?
- Host multiple CSV files on a cloud provider to request through both HTTP and HTTPS (possibly even decline data hosted over HTTP).
- Create a suite of "bad files" containing errors like (1) empty file, (2) not CSV, (3) missing data, (4) malformed CSV, (5) improper headers, etc. Could be created through fuzzing.
- Have a scalable system to spin up multiple machines.
- Have distributed data to test the machines ability to handle latency.
- Vary traffic loads
- Vary file sizes
