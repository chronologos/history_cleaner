package fixer

import (
	"fmt"
	"os"
	"testing"

	"github.com/r3labs/diff"
)

type stubStringWriter struct{}

func (*stubStringWriter) WriteString(_ string) (int, error) {
	return 0, nil
}

func TestFixer(t *testing.T) {
	testCases := []struct {
		name string
		inHistory   []string
		wantHistory []string

	}{
		{
		name:        "empty",
		inHistory:   []string{},
		wantHistory: []string{},
	},
		{
			name: "basic 1",
			inHistory:   []string{"#1", "ls", "#2", "ls", "#2", "cd"},
			wantHistory: []string{"#1", "ls", "#2", "cd"},
		},
		{
			name: "basic 2",
			inHistory:   []string{"#1", "ls -al", "#2", "cd", "#3", "ls -l", "#3", "ls -al", "#4", "cd", "#5"},
			wantHistory: []string{"#1", "ls -al", "#2", "cd", "#3", "ls -l"},
		},
		{
			name: "remove temp tag",
			inHistory:   []string{"#1", fmt.Sprintf("ls -al %s", tempTag),"#2", "cd",},
			wantHistory: []string{"#2", "cd"},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test %s", tc.name), func(t *testing.T) {
			bob := New(tc.inHistory, os.Stdout)
			//bob := New(tc.inHistory, &stubStringWriter{})
			gotCleanedHistory := bob.Fix()
			if changelog, _ := diff.Diff(gotCleanedHistory, tc.wantHistory); len(changelog)!=0 {
				t.Errorf("history != cleanedHistory, diff=%v", changelog)
			}
		})
	}
}
