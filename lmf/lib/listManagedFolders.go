package lib

import (
	"path/filepath"
	"regexp"
)

// GetScriptPath joins the autoGenFolder with cleaned versions of location1 and location2.
// In Go, filepath.Join is the equivalent of Java's Path.resolve().
func GetScriptPath(autoGenFolder string, location1 string, location2 string) string {
	return filepath.Join(
		autoGenFolder,
		GetCleanLocationName(location1),
		GetCleanLocationName(location2),
	)
}

// GetCleanLocationName replaces special characters with underscores and trims trailing ones.
func GetCleanLocationName(location string) string {
	// Replaces one or more of ':', '\', '/', or ' ' with "_"
	reLegal := regexp.MustCompile(`[:\\/ ]+`)
	allLegal := reLegal.ReplaceAllString(location, "_")

	// Replaces one or more trailing underscores at the end of the string
	reTrailing := regexp.MustCompile(`_+$`)
	return reTrailing.ReplaceAllString(allLegal, "")
}
