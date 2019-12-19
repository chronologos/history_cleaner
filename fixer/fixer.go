package fixer

import "errors"

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
