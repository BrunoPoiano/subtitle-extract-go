package types

import "os"

// SubtitleJob represents a task to extract subtitles from a video file.
// It contains the file location and information about the file itself.
type SubtitleJob struct {
	Location string      // Directory path where the video file is located
	Item     os.DirEntry // File information including name and metadata
}
