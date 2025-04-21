package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// rootFolder defines the base directory where videos will be searched.
// Production and development paths are provided - uncomment as needed.
var rootFolder = "/app/subextract/videos"

// videoExtensions is a map of supported video file extensions for subtitle extraction.
// Files with these extensions will be processed by the application.
var videoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
}

// SubtitleJob represents a task to extract subtitles from a video file.
// It contains the file location and information about the file itself.
type SubtitleJob struct {
	location string      // Directory path where the video file is located
	item     os.DirEntry // File information including name and metadata
}

// runWorkerPool creates and manages a pool of worker goroutines to process subtitle extraction jobs.
// Parameters:
//   - jobs: channel of subtitle extraction tasks
//   - workers: number of concurrent workers to spawn
//   - wg: wait group to track job completion
func runWorkerPool(jobs <-chan SubtitleJob, workers int, wg *sync.WaitGroup) {
	for range workers {
		go func() {
			for job := range jobs {
				runningEmbedSubtitleCheck(job.location, job.item)
				wg.Done()
			}
		}()
	}
}

// main is the entry point of the program that orchestrates the subtitle extraction process.
// It sets up worker pools for parallel processing and waits for all jobs to complete.
func main() {
	var wg sync.WaitGroup
	cpus_available := runtime.NumCPU()

	if cpus_available > 1 {
		cpus_available--
	}

	println("using", cpus_available, "workers")

	jobs := make(chan SubtitleJob)

	// Initialize worker pool to process subtitle extraction jobs concurrently
	runWorkerPool(jobs, cpus_available, &wg)

	// Search for video files and enqueue extraction jobs
	getItemsFromFolder(rootFolder, jobs, &wg)

	close(jobs) // Signal that no more jobs will be added

	wg.Wait() // Wait for all extraction jobs to complete
}

// getItemsFromFolder recursively processes all items in the given directory,
// looking for video files to extract subtitles from. It adds subtitle extraction
// jobs to the job queue when it finds video files without existing subtitles.
// Parameters:
//   - location: directory path to search in
//   - jobs: channel to send subtitle extraction tasks
//   - wg: wait group to track job submission
func getItemsFromFolder(location string, jobs chan<- SubtitleJob, wg *sync.WaitGroup) {
	items, err := os.ReadDir(location)
	if err != nil {
		fmt.Printf("Error accessing folder %s: %v\n", location, err)
		return
	}

	for _, item := range items {

		if item.IsDir() {
			newLoc := newLocation(location, item.Name())
			getItemsFromFolder(newLoc, jobs, wg) // Recursively process subdirectories
			continue
		}

		ext := extractExtention(item.Name())

		if videoExtensions[ext] {

			// Create full path for potential English subtitle file
			subName := strings.TrimSuffix(item.Name(), ext)
			newPath := newLocation(location, subName)
			newName := newSrtName(newPath, "eng")

			// Check if subtitle file already exists to avoid redundant extraction
			_, err := os.Stat(newName)
			if err != nil {
				wg.Add(1)                           // Register new job with wait group
				jobs <- SubtitleJob{location, item} // Queue the job for processing
			} else {
				println("subtitles already extracted for:", item.Name())
			}
		}
	}
}

// runningEmbedSubtitleCheck examines a video file for embedded subtitles using ffmpeg.
// If subtitles are found, it identifies each subtitle stream and language code,
// then initiates the extraction process for each detected subtitle track.
// Parameters:
//   - location: directory path containing the video file
//   - item: file information of the video to check
func runningEmbedSubtitleCheck(location string, item os.DirEntry) {

	fullPath := newLocation(location, item.Name())

	// Run ffmpeg to get file metadata including subtitle information
	cmd := exec.Command("ffmpeg", "-i", fullPath)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "Stream") {
		println("No subtitles found for |", item.Name())
		return
	}

	println("Extracting from |", item.Name())

	// Regex to detect subtitle streams with language codes
	reStream := regexp.MustCompile(`Stream #\d+:\d+\((\w{3})\): Subtitle`)
	matches := reStream.FindAllStringSubmatch(string(output), -1)

	if len(matches) == 0 {
		println("No subtitle stream detected for |", item.Name())
		return
	}

	// Regex to extract stream index and language code
	reIndex := regexp.MustCompile(`Stream #(\d+:\d+)\((\w{3})\): Subtitle`)
	languages := strings.Split(string(output), "\n")

	// Process each line of ffmpeg output looking for subtitle streams
	for _, language := range languages {

		if !strings.Contains(language, "Subtitle") {
			continue
		}

		match := reIndex.FindStringSubmatch(language)

		if len(match) != 3 {
			continue
		}

		subtitleIndex := match[1]    // Stream index (e.g., "0:2")
		subtitleLanguage := match[2] // Language code (e.g., "eng")

		runningExtractSubtitle(location, item.Name(), subtitleIndex, subtitleLanguage)
	}

}

// runningExtractSubtitle extracts embedded subtitles from a video file
// using ffmpeg with the specified subtitle track and language code.
// It saves the subtitle as an SRT file with appropriate naming.
// Parameters:
//   - location: directory path containing the video file
//   - name: filename of the video
//   - subtitleIndex: ffmpeg stream index for the subtitle track (e.g. "0:2")
//   - subtitleLanguage: three-letter language code (e.g. "eng")
func runningExtractSubtitle(location, name, subtitleIndex, subtitleLanguage string) {

	fullPath := newLocation(location, name)

	newName := strings.TrimSuffix(fullPath, extractExtention(name))
	fullName := newSrtName(newName, subtitleLanguage)

	// Extract the subtitle track and convert to SRT format
	cmd := exec.Command("ffmpeg", "-i", fullPath, "-map", subtitleIndex, "-c:s", "srt", fullName)

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	println("extracted |", subtitleLanguage, "from", name)

}

// newSrtName creates a subtitle filename by combining a base path and language code.
// For example, "/path/movie" + "eng" becomes "/path/movie.eng.srt"
// Parameters:
//   - value1: base path and filename without extension
//   - value2: language code to append
//
// Returns:
//   - complete .srt filename with language code
func newSrtName(value1, value2 string) string {
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
func newLocation(value1, value2 string) string {
	return filepath.Join(value1, value2)
}

// extractExtention returns the lowercase file extension of a filename.
// For example, "movie.MP4" returns ".mp4" (normalized to lowercase)
// Parameters:
//   - value: filename including extension
//
// Returns:
//   - lowercase extension with dot prefix
func extractExtention(value string) string {
	return strings.ToLower(filepath.Ext(value))
}
