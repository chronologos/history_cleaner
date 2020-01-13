package main

import (
	"bufio"
	"fmt"
	"github.com/chronologos/history_cleaner/fixer"
	"log"
	"os"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

var onlyOnce = []string{}

const (
	//PATH_TO_HISTORY = "/usr/local/google/home/iantay/.bash_history"
	PATH_TO_HISTORY = "/Users/iantay/.bash_history"
	//PATH_TO_LOGS = "/usr/local/google/home/iantay/usage.log"
	PATH_TO_LOGS = "/Users/iantay/usage.log"
)

func main() {
	inFile, err := os.Open(PATH_TO_HISTORY)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	logFile, err := os.OpenFile(PATH_TO_LOGS, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logWriter := bufio.NewWriter(logFile)
	defer logWriter.Flush()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	var res []string
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	histFixer := fixer.New(res, logWriter)
	cleanedHistory := histFixer.Fix()

	if err := os.Remove(PATH_TO_HISTORY); err != nil {
		log.Fatal(err)
	}
	outFile, err := os.OpenFile(PATH_TO_HISTORY, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	w := bufio.NewWriter(outFile)
	for _, l := range cleanedHistory {
		w.WriteString(l)
		w.WriteString("\n")
	}
	w.Flush()
	fmt.Println("done!")
}