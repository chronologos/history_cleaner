package fixer

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
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
			name: "basic 1",
			inHistory:   []string{"#1", "ls", "#2", "ls", "#2", "cd"},
			wantHistory: []string{"#1", "ls", "#2", "cd"},
		},
		{
			name: "basic 2",
			inHistory:   []string{"#1", "ls -al", "#2", "cd", "#3", "ls -l", "#3", "ls -al", "#4", "cd", "#5"},
			wantHistory: []string{"#1", "ls -al", "#2", "cd", "#3", "ls -l"},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test %s", tc.name), func(t *testing.T) {
			bob := New(tc.inHistory, os.Stdout)
			//bob := New(tc.inHistory, &stubStringWriter{})
			gotCleanedHistory := bob.Fix()
			if !reflect.DeepEqual(gotCleanedHistory, tc.wantHistory) {
				t.Errorf("history != cleanedHistory, got=%v, want=%v", strings.Join(gotCleanedHistory, ","), strings.Join(tc.wantHistory, ","))
			}
		})
	}
}
