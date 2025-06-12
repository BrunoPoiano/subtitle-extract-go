package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// newSrtName creates a subtitle filename by combining a base path and language code.
// For example, "/path/movie" + "eng" becomes "/path/movie.eng.srt"
// Parameters:
//   - value1: base path and filename without extension
//   - value2: language code to append
//
// Returns:
//   - complete .srt filename with language code
func NewSrtName(value1, value2 string) string {
	return fmt.Sprintf("%s.%s.srt", value1, value2)
}

// newLocation joins two path components to create a new file path.
// This is a wrapper around filepath.Join for cleaner code.
// Parameters:
//   - value1: base directory path
//   - value2: filename or subdirectory to append
//
// Returns:
//   - combined path that is platform-appropriate
func NewLocation(value1, value2 string) string {
	return filepath.Join(value1, value2)
}

// extractExtention returns the lowercase file extension of a filename.
// For example, "movie.MP4" returns ".mp4" (normalized to lowercase)
// Parameters:
//   - value: filename including extension
//
// Returns:
//   - lowercase extension with dot prefix
func ExtractExtention(value string) string {
	return strings.ToLower(filepath.Ext(value))
}

// generateLogs formats and prints log messages with a timestamp prefix.
// It takes a log message as input, adds the current date and time in the format "DD/MM/YYYY HH:MM:SS",
// and outputs the combined string to standard output.
// Parameters:
//   - value: the log message to be printed
func GenerateLogs(value string) {
	now := time.Now()
	timeFormated := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	println(timeFormated, "|", value)
}
