package minishell

import (
	"os"
	"reflect"
	"testing"
)

func TestCd(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "existing dir",
			input:    "..",
			expected: "/home/ibragim/myProjects/wb-tech/go-minishell/internal",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cd(tc.input)
			wd, _ := os.Getwd()
			if !reflect.DeepEqual(wd, tc.expected) {
				t.Errorf("cd() = %v, want: %v", wd, tc.expected)
			}
		})
	}
}

func TestPwd(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "get pwd",
			expected: "/home/ibragim/myProjects/wb-tech/go-minishell/internal/minishell",
		},
	}

	for _, tc := range testCases {
		os.Chdir("/home/ibragim/myProjects/wb-tech/go-minishell/internal/minishell")
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := pwd()
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("pwd() = %v, want: %v", res, tc.expected)
			}
		})
	}
}
