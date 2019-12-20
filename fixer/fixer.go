package fixer

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type Scanner interface {
	Scan() bool
	Text() string
}

type Config struct {
	OnlyOnce []string
}

type Bob struct {
	history []string
	logger  io.StringWriter
	config  Config
}

func New(history []string, logger io.StringWriter) *Bob {
	return &Bob{
		history: history,
		logger:  logger,
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

type iterator struct {
	iterable []string
	length int
	index int
}

func NewIterator(in []string) *iterator {
	fmt.Println("cons")
	copyIn := in[:]
	return &iterator{
		iterable:copyIn,
		length: len(copyIn),
	}
}

func (iter *iterator) GetNext() (int, string, error) {
	fmt.Printf("%+v\n",iter)
	if iter.index >= iter.length {
		return 0, "", io.EOF
	}
	oldIndex := iter.index
	iter.index+=1
	return oldIndex, iter.iterable[oldIndex], nil
}


func (b *Bob) Fix() []string {
	// map of command to timestamp
	seen := make(map[string]*timestamp)
	var out []string
	shouldBeTimestamp := true
	var curTimestamp *timestamp
	var allCommands int
	onlyOnceCounter := make([]int, len(b.config.OnlyOnce))
	currentTime := asTimestamp(fmt.Sprintf("#%v", time.Now().Unix()))
	iter := NewIterator(b.history)
	
	var err error
	var i int
	var l string
OUTER:
	for err == nil {
	i, l, err = iter.GetNext()
	//for i, l := range b.history {
		fmt.Println("loop")
		if shouldBeTimestamp {
			shouldBeTimestamp = false
			if err := IsValidTimestamp(l); err == nil {
				b.logger.WriteString("got ts\n")
				curTimestamp = asTimestamp(l)
				continue OUTER
			} else {
				b.logger.WriteString(fmt.Sprintf("on date: %v: failed to parse history file at line %d = %s. got err=%v\n", time.Now(), i, l, err))
				curTimestamp = currentTime
			}
		}
		if !shouldBeTimestamp {
			shouldBeTimestamp = true
			if err := IsValidCommand(l); err != nil {
				b.logger.WriteString(fmt.Sprintf("on date: %v: failed to parse history file at line %d. got timestamp %s, expected command\n", time.Now(), i, l))
				continue OUTER
			}
			allCommands++
			strippedLine := strings.TrimSpace(l)
			if _, ok := seen[strippedLine]; !ok {
				seen[strippedLine] = curTimestamp
				for i, substr := range b.config.OnlyOnce {
					if strings.Contains(strippedLine, substr) && onlyOnceCounter[i] != 0 {
						continue OUTER
					}
					onlyOnceCounter[i] += 1
				}
				out = append(out, curTimestamp.String(), strippedLine)
			}
		}

	}
	fmt.Printf("on date: %v: saw %d commands, %d unique, %d removed\n", time.Now(), allCommands, len(seen), allCommands-len(seen))
	b.logger.WriteString(fmt.Sprintf("on date: %v: saw %d commands, %d unique, %d removed\n", time.Now(), allCommands, len(seen), allCommands-len(seen)))
	return out
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
