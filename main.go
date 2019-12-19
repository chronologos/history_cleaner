package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	fixer "github.com/chronologos/historyfixer/fixer"
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
	PATH_TO_HISTORY = "/usr/local/google/home/iantay/.bash_history"
)

func main() {
	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	inFile, _ := os.Open(PATH_TO_HISTORY)
	defer inFile.Close()

	logFile, err := os.OpenFile("/usr/local/google/home/iantay/usage.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	logwriter := bufio.NewWriter(logFile)
	defer logwriter.Flush()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	res := []string{}
	for scanner.Scan() {
		res = append(res, scanner.Text())

	}
	// map of command to timestamp
	seen := make(map[string]string)
	var out []string
	shouldBeTimestamp := true
	var curTimestamp string
	var allCommands int
	onlyOnceCounter := make([]int, len(onlyOnce))
	now := time.Now()
	secs := fmt.Sprintf("#%v", now.Unix())
OUTER:
	for i, l := range res {
		if shouldBeTimestamp {
			shouldBeTimestamp = false
			if err := fixer.IsValidTimestamp(l); err == nil {
				curTimestamp = l
				continue
			} else {
				logwriter.WriteString(fmt.Sprintf("on date: %v: failed to parse history file at line %d. got err=%v\n", time.Now(), i, err))
				curTimestamp = secs
			}
		}
		if !shouldBeTimestamp {
			shouldBeTimestamp = true
			if err := fixer.IsValidCommand(l); err != nil {
				logwriter.WriteString(fmt.Sprintf("on date: %v: failed to parse history file at line %d. got timestamp %s, expected command\n", time.Now(), i, l))
				continue OUTER
			}
			allCommands++
			strippedLine := strings.TrimSpace(l)
			if _, ok := seen[strippedLine]; !ok {
				seen[strippedLine] = curTimestamp
				for i, substr := range onlyOnce {
					if strings.Contains(strippedLine, substr) && onlyOnceCounter[i] != 0 {
						continue OUTER
					}
					onlyOnceCounter[i] += 1
				}
				out = append(out, curTimestamp, strippedLine)
			}
		}

	}
	fmt.Printf("on date: %v: saw %d commands, %d unique, %d removed\n", time.Now(), allCommands, len(seen), allCommands-len(seen))
	logwriter.WriteString(fmt.Sprintf("on date: %v: saw %d commands, %d unique, %d removed\n", time.Now(), allCommands, len(seen), allCommands-len(seen)))

	if err := os.Remove(PATH_TO_HISTORY); err != nil {
		log.Fatal(err)
	}
	outFile, err := os.OpenFile(PATH_TO_HISTORY, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	w := bufio.NewWriter(outFile)
	for _, l := range out {
		w.WriteString(l)
		w.WriteString("\n")
	}
	w.Flush()
	fmt.Println("done!")
}

func isValidTimestamp(t string) error {
	if len(t) == 0 {
		return errors.New("invalid timestamp: no data")
	}
	if t[0] != '#' {
		return errors.New("invalid timestamp: does not start with #")
	}
	return nil
}

func isValidCommand(c string) error {
	if len(c) == 0 {
		return errors.New("invalid command: no data")
	}
	if c[0] == '#' {
		return errors.New("invalid command: starts with #")
	}
	return nil
}
