package subtitleService

import (
	"fmt"
	"main/config"
	"main/types"
	"main/utils"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// getItemsFromFolder recursively processes all items in the given directory,
// looking for video files to extract subtitles from. It adds subtitle extraction
// jobs to the job queue when it finds video files without existing subtitles.
// Parameters:
//   - location: directory path to search in
//   - jobs: channel to send subtitle extraction tasks
//   - wg: wait group to track job submission
func GetItemsFromFolder(location string, jobs chan<- types.SubtitleJob, wg *sync.WaitGroup) {
	items, err := os.ReadDir(location)
	if err != nil {
		utils.GenerateLogs(fmt.Sprintf("ERROR | accessing folder %s: %v\n", location, err))
		return
	}

	for _, item := range items {

		if item.IsDir() {
			newLoc := utils.NewLocation(location, item.Name())
			GetItemsFromFolder(newLoc, jobs, wg) // Recursively process subdirectories
			continue
		}

		ext := utils.ExtractExtention(item.Name())

		if config.VideoExtensions[ext] {

			// Create full path for potential English subtitle file
			subName := strings.TrimSuffix(item.Name(), ext)
			newPath := utils.NewLocation(location, subName)
			newName := utils.NewSrtName(newPath, "eng")

			// Check if subtitle file already exists to avoid redundant extraction
			_, err := os.Stat(newName)
			if err != nil {
				wg.Add(1)                                                 // Register new job with wait group
				jobs <- types.SubtitleJob{Location: location, Item: item} // Queue the job for processing
			} else {
				utils.GenerateLogs(fmt.Sprintf("INFO | subtitles already extracted | %s", item.Name()))
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
func RunningEmbedSubtitleCheck(location string, item os.DirEntry) {

	fullPath := utils.NewLocation(location, item.Name())

	// Run ffmpeg to get file metadata including subtitle information
	cmd := exec.Command("ffmpeg", "-i", fullPath)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "Stream") {
		utils.GenerateLogs(fmt.Sprintf("INFO | No subtitles found | %s", item.Name()))
		return
	}

	utils.GenerateLogs(fmt.Sprintf("PROCESS | Extracting from | %s", item.Name()))

	// Regex to detect subtitle streams with language codes
	reStream := regexp.MustCompile(`Stream #\d+:\d+\((\w{3})\): Subtitle`)
	matches := reStream.FindAllStringSubmatch(string(output), -1)

	if len(matches) == 0 {
		utils.GenerateLogs(fmt.Sprintf("INFO | No subtitle detected | %s", item.Name()))
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
	runSubtitleSync(location, item.Name())
}

// runSubtitleSync synchronizes subtitles using the alass-cli tool.
// It takes the location and name of a video file as input.  It first runs alass-cli to generate a "fixed" subtitle file. Then, it iterates through all subtitle files in the directory, using alass-cli to synchronize each with the fixed subtitle file. Finally, it removes the temporary "fixed" subtitle file.
func runSubtitleSync(location string, name string) {

	items, err := os.ReadDir(location)
	if err != nil {
		utils.GenerateLogs(fmt.Sprintf("ERROR | %s ", err))
		return
	}

	fullPath := utils.NewLocation(location, name)
	newName := strings.TrimSuffix(fullPath, utils.ExtractExtention(name))
	fullEngName := utils.NewSrtName(newName, "eng")
	fullFixName := utils.NewSrtName(newName, "fix")

	//generate the fix subtitle based on the eng subtitle
	cmd := exec.Command("./alass-cli", fullPath, fullEngName, fullFixName)
	_, err = cmd.CombinedOutput()
	if err != nil {
		utils.GenerateLogs(fmt.Sprintf("ERROR | %s ", err))
		return
	}

	for _, item := range items {

		ext := utils.ExtractExtention(item.Name())
		if config.SubtitleExtensions[ext] {

			itemPath := utils.NewLocation(location, item.Name())

			//use the fix subtitle as reference to synch all the subtitles available for this item
			cmd := exec.Command("./alass-cli", fullPath, fullFixName, itemPath)
			_, err = cmd.CombinedOutput()
			if err != nil {
				utils.GenerateLogs(fmt.Sprintf("ERROR | %s ", err))
				continue
			}
			utils.GenerateLogs(fmt.Sprintf("INFO | Synched | %s", item.Name()))
		}
	}

	// remove the fix subtitle
	cmd = exec.Command("rm", fullFixName)
	_, err = cmd.CombinedOutput()
	if err != nil {
		utils.GenerateLogs(fmt.Sprintf("ERROR | %s ", err))
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

	fullPath := utils.NewLocation(location, name)

	newName := strings.TrimSuffix(fullPath, utils.ExtractExtention(name))
	fullName := utils.NewSrtName(newName, subtitleLanguage)

	// Extract the subtitle track and convert to SRT format
	cmd := exec.Command("ffmpeg", "-i", fullPath, "-map", subtitleIndex, "-c:s", "srt", fullName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		utils.GenerateLogs(fmt.Sprintf("ERROR | %s ", err))
		return
	}

	utils.GenerateLogs(fmt.Sprintf("INFO | extracted | %s | %s", subtitleLanguage, name))

}
