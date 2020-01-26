package fixer

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	tempTag = "#TEMP"
)
type Config struct {
	OnlyOnce []string
}

type Fixer struct {
	history []string
	logger  io.StringWriter
	config  Config

	// Needed to enable main functionality.
	commandTime *timestamp
	currentTime *timestamp
	expectingTimestamp bool
	seen map[string]*timestamp
	onlyOnceCount []int
	commandCount int

	stringsToRemove []string
}

func New(history []string, logger io.StringWriter) *Fixer {
	return &Fixer{
		history: history,
		logger:  logger,
		currentTime: asTimestamp(fmt.Sprintf("#%v", time.Now().Unix())),
		seen: make(map[string]*timestamp),
		expectingTimestamp: true,
		stringsToRemove: []string{tempTag},
	}
}

type timestamp struct {
	t string
}

func asTimestamp(s string) *timestamp {
	return &timestamp{s}
}

func (t *timestamp) String() string {
	return t.t
}

func (fixer *Fixer) Fix() []string {
	var out []string

	for i, l := range fixer.history {
		out = append(out, fixer.processOneLine(i, l)...)
	}

	fmt.Printf("on date: %v: saw %d commands, %d unique, %d removed\n", time.Now(), fixer.commandCount, len(fixer.seen), fixer.commandCount-len(fixer.seen))
	fixer.logger.WriteString(fmt.Sprintf("on date: %v: saw %d commands, %d unique, %d removed\n", time.Now(), fixer.commandCount, len(fixer.seen), fixer.commandCount-len(fixer.seen)))
	return out
}

func (fixer *Fixer) processOneLine(i int, l string) []string {
	if fixer.expectingTimestamp {
		fixer.expectingTimestamp = false
		if err := IsValidTimestamp(l); err == nil {
			//fixer.logger.WriteString("got ts\n")
			fixer.commandTime = asTimestamp(l)
		} else {
			fixer.logger.WriteString(fmt.Sprintf("on date: %v: error in history file at line %d = %s. got err=%v\n", time.Now(), i, l, err))
			fixer.commandTime = fixer.currentTime
		}
		return []string{}
	}
	fixer.expectingTimestamp = true
		if err := IsValidCommand(l); err != nil {
			fixer.logger.WriteString(fmt.Sprintf("on date: %v: failed to parse history file at line %d. got timestamp %s, expected command\n", time.Now(), i, l))
			return []string{}
		}
		fixer.commandCount++
		strippedLine := strings.TrimSpace(l)
		for _, dontWant := range fixer.stringsToRemove {
			if strings.Contains(strippedLine, dontWant){
				return []string{}
			}
		}
		if _, ok := fixer.seen[strippedLine]; !ok {
			fixer.seen[strippedLine] = fixer.commandTime
			for i, substr := range fixer.config.OnlyOnce {
				if strings.Contains(strippedLine, substr) && fixer.onlyOnceCount[i] != 0 {
					return []string{}
				}
				fixer.onlyOnceCount[i] += 1
			}
			if fixer.commandTime != nil {
				return []string{fixer.commandTime.String(), strippedLine}
			}
		}
		return []string{}
}

func IsValidTimestamp(t string) error {
	if len(t) == 0 {
		return errors.New("invalid timestamp: no data")
	}
	if t[0] != '#' {
		return errors.New("invalid timestamp: does not start with #")
	}
	return nil
}

func IsValidCommand(c string) error {
	if len(c) == 0 {
		return errors.New("invalid command: no data")
	}
	if c[0] == '#' {
		return errors.New("invalid command: starts with #")
	}
	return nil
}
