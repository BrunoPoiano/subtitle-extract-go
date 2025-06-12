package config

import (
	"main/utils"
	"regexp"
	"strings"
)

// rootFolder defines the base directory where videos will be searched.
var RootFolder = "/app/subextract/videos"

// videoExtensions is a map of supported video file extensions for subtitle extraction.
// Files with these extensions will be processed by the application.
var VideoExtensions = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".wmv":  true,
	".flv":  true,
	".webm": true,
}

// SubtitleExtensions defines the supported subtitle file extensions.
var SubtitleExtensions = map[string]bool{
	".srt": true,
}

// SubtitlesToExtract retrieves a map of subtitle languages to extract from the environment variable SUBTITLES_TO_EXTRACT.
// The environment variable should contain a comma-separated list of language codes enclosed in square brackets, e.g., "[por,eng]".
// If the environment variable is empty or contains "[]", an empty map is returned.
func SubtitlesToExtract() map[string]bool {
	subtitle := make(map[string]bool)
	raw := utils.ReturnEnvVariable("SUBTITLES_TO_EXTRACT", "[]")

	if raw == "" || raw == "[]" {
		return subtitle
	}

	raw = regexp.MustCompile(`[\[\]]`).ReplaceAllString(raw, "")

	for _, item := range strings.Split(raw, ",") {
		subtitle[strings.TrimSpace(item)] = true
	}
	return subtitle
}

// DefaultSub retrieves the default subtitle language code from the environment variable DEFAULT_SUB.
// If the environment variable is not set, it defaults to "eng".
func DefaultSub() string {
	return utils.ReturnEnvVariable("DEFAULT_SUB", "eng")
}
