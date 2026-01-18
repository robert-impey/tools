package lib

import (
	"path/filepath"
	"testing"
)

func TestGetCleanLocationName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"C:\\Data", "C_Data"},
		{"D:/Archive/", "D_Archive"},
		{"E:\\My Documents", "E_My_Documents"},
	}

	for _, tt := range tests {
		result := GetCleanLocationName(tt.input)
		if result != tt.expected {
			t.Errorf("GetCleanLocationName(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetScriptPath(t *testing.T) {
	autoGen := "C:\\AutoGen"
	// filepath.Join handles platform-specific separators automatically
	expected := filepath.Join(autoGen, "C_Data", "D_Backup")

	result := GetScriptPath(autoGen, "C:\\Data", "D:\\Backup")

	if result != expected {
		t.Errorf("GetScriptPath() = %q; want %q", result, expected)
	}
}
