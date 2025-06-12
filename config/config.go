package config

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

var SubtitleExtensions = map[string]bool{
	".srt": true,
}
