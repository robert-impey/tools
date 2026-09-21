package internal

import (
	"path/filepath"
	"testing"
)

func TestRsyncFileHasCopies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileName string
		want     bool
	}{
		{
			name:     "receiving log with copied files",
			fileName: "receiving-copies.log",
			want:     true,
		},
		{
			name:     "sending log with copied files",
			fileName: "sending-copies.log",
			want:     true,
		},
		{
			name:     "log without copied files",
			fileName: "empty.log",
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePath := filepath.Join("test_data", "rsync", test.fileName)

			hasCopies, err := RsyncFileHasCopies(filePath)
			if err != nil {
				t.Fatalf("RsyncFileHasCopies error: %v", err)
			}
			if hasCopies != test.want {
				t.Fatalf("expected RsyncFileHasCopies to return %t for %s", test.want, test.fileName)
			}
		})
	}
}
