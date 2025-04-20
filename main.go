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
var rootFolder = "/app/subextract/videos"
//var rootFolder = "/home/brunopoiano/Documents/Pessoal/sub-extract-go/videos"

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
			println("----------------------------------------")
		}
	}
}

// runningEmbedSubtitleCheck checks if a video file has embedded subtitles
// using ffmpeg and initiates extraction if subtitles are found
func runningEmbedSubtitleCheck(location string, item os.DirEntry) {

	fullPath := newLocation(location, item.Name())

	cmd := exec.Command("ffmpeg", "-i", fullPath)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "Stream") {
		println("No subtitles found for:", item.Name())
		return
	}

	println("Extracting subtitles for:", item.Name())

	reStream := regexp.MustCompile(`Stream #\d+:\d+\((\w{3})\): Subtitle`)
	matches := reStream.FindAllStringSubmatch(string(output), -1)

	if len(matches) == 0 {
		println("No subtitle stream detected")
		return
	}

	reIndex := regexp.MustCompile(`Stream #(\d+:\d+)\((\w{3})\): Subtitle`)
	languages := strings.Split(string(output), "\n")

	for _, language := range languages {

		if !strings.Contains(language, "Subtitle") {
			continue
		}

		match := reIndex.FindStringSubmatch(language)

		if len(match) != 3 {
			continue
		}

		subtitleIndex := match[1]
		subtitleLanguage := match[2]

		runningExtractSubtitle(location, item.Name(), subtitleIndex, subtitleLanguage)
	}

}

// runningExtractSubtitle extracts embedded subtitles from a video file
// using ffmpeg with the specified subtitle track and language
func runningExtractSubtitle(location, name, subtitleIndex, subtitleLanguage string) {

	fullPath := newLocation(location, name)

	newName := strings.TrimSuffix(fullPath, extractExtention(name))
	fullName := newSrtName(newName, subtitleLanguage)

	cmd := exec.Command("ffmpeg", "-i", fullPath, "-map", subtitleIndex, "-c:s", "srt", fullName)

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	println(subtitleLanguage, "extracted")

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
