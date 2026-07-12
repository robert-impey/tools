package internal

import (
	"testing"
)

func TestCleanFolderPathForLogName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/folder", "path_to_folder"},
		{"C:\\path\\to\\folder", "C_path_to_folder"},
		{"user@host:/path/to/folder", "user__host__path_to_folder"},
		{"/Users/robert/Documents/Project:Name", "Users_robert_Documents_ProjectName"},
		{"Z:\\", "Z"},
		{"Z:", "Z"},
	}

	for _, test := range tests {
		result := CleanFolderPathForLogName(test.input)
		if result != test.expected {
			t.Errorf("CleanFolderPathForLogName(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
