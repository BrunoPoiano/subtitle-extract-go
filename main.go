package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// rootFolder defines the base directory where videos will be searched
var rootFolder = "/subextract/videos"

// videoExtensions is a map of supported video file extensions
var videoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
}

// main is the entry point of the program, starts the subtitle extraction process
func main() {
	os.Setenv("PROD", "PRODUCTION")
	getItemsFromFolder(rootFolder)
}

// getItemsFromFolder recursively processes all items in the given directory
// looking for video files to extract subtitles from
func getItemsFromFolder(location string) {

	items, err := os.ReadDir(location)
	if err != nil {
		fmt.Printf("Error accessing folder %s: %v\n", location, err)
		return
	}

	for _, item := range items {

		if item.IsDir() {
			newLocation := newLocation(location, item.Name())
			getItemsFromFolder(newLocation)
			continue
		}

		ext := extractExtention(item.Name())

		if videoExtensions[ext] {

			// creating fullPath of eng subtitle
			subName := strings.TrimSuffix(item.Name(), ext)
			newPath := newLocation(location, subName)
			newName := newSrtName(newPath, "eng")

			//checking if file already exists
			_, err := os.Stat(newName)
			if err != nil {
				runningEmbedSubtitleCheck(location, item)
			} else {
				println("subtitles already extracted for:", item.Name())
			}
		}

	}
}

// runningEmbedSubtitleCheck checks if a video file has embedded subtitles
// using ffmpeg and initiates extraction if subtitles are found
func runningEmbedSubtitleCheck(location string, item os.DirEntry) {

	name := item.Name()
	command := fmt.Sprintf("ffmpeg -i '%s/%s' 2>&1 | grep Stream.*Subtitle", location, name)

	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		println("No subtitles found for:", item.Name())
		return
	}

	println("Extracting subtitles for:", item.Name())

	items := strings.Split(string(output), "\n")

	for _, item := range items {
		re := regexp.MustCompile(`\d+:\d+\([a-zA-Z0-9]{3}\)`)
		matches := re.FindAllString(item, -1)

		if len(matches) == 0 {
			continue
		}
		runningExtractSubtitle(location, name, matches[0])
	}

}

// runningExtractSubtitle extracts embedded subtitles from a video file
// using ffmpeg with the specified subtitle track and language
func runningExtractSubtitle(location, name, subLanguage string) {

	re := regexp.MustCompile(`[()]`)

	subDiv := strings.Split(subLanguage, "(")
	subNum := string(subDiv[0])
	subLen := re.ReplaceAllString(subDiv[1], "")

	fullPath := newLocation(location, name)

	newName := strings.TrimSuffix(fullPath, extractExtention(name))
	fullName := newSrtName(newName, subLen)

	command := fmt.Sprintf("ffmpeg -i '%s' -map %s -c:s srt '%s'", fullPath, subNum, fullName)

	cmd := exec.Command("bash", "-c", command)

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

}

// newSrtName creates a subtitle filename by combining a base path and language code
func newSrtName(value1, value2 string) string {
	return fmt.Sprintf("%s.%s.srt", value1, value2)
}

// newLocation joins two path components to create a new file path
func newLocation(value1, value2 string) string {
	return filepath.Join(value1, value2)
}

// extractExtention returns the lowercase file extension of a filename
func extractExtention(value string) string {
	return strings.ToLower(filepath.Ext(value))
}
