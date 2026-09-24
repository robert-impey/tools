package internal

import (
	"path/filepath"
	"testing"
)

func TestRsyncFileHasCopies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		fileName   string
		wantRsync  bool
		wantCopies bool
	}{
		{
			name:       "receiving log with copied files",
			fileName:   "receiving-copies.log",
			wantRsync:  true,
			wantCopies: true,
		},
		{
			name:       "sending log with copied files",
			fileName:   "sending-copies.log",
			wantRsync:  true,
			wantCopies: true,
		},
		{
			name:       "log without copied files",
			fileName:   "empty.log",
			wantRsync:  true,
			wantCopies: false,
		},
		{
			name:       "unrelated log file",
			fileName:   "gen-script.log",
			wantRsync:  false,
			wantCopies: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			filePath := filepath.Join("test_data", "rsync", test.fileName)

			isRsyncLog, hasCopies, err := RsyncFileHasCopies(filePath)
			if err != nil {
				t.Fatalf("RsyncFileHasCopies error: %v", err)
			}
			if isRsyncLog != test.wantRsync {
				t.Fatalf("expected isRsyncLog to be %t for %s", test.wantRsync, test.fileName)
			}
			if hasCopies != test.wantCopies {
				t.Fatalf("expected hasCopies to be %t for %s", test.wantCopies, test.fileName)
			}
		})
	}
}
