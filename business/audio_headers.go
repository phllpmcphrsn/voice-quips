package business

import "strings"

// Constants for the content-type headers for audio files
const (
	AudioHeader = "audio/"
	MP3Header   = AudioHeader + "mpeg" // audio/mpeg
	MP4Header   = AudioHeader + "mp4"  // audio/mp4
	WAVHeader   = AudioHeader + "wav"  // audio/wav
	OGGHeader   = AudioHeader + "ogg"  // audio/ogg
	FLACHeader  = AudioHeader + "flac" // audio/flac
)

// GetContentType returns the content-type header for a given audio file extension
func GetContentType(ext string) string {
	// convert the extension to lower case
	ext = strings.ToLower(ext)

	// use a switch statement with the AudioHeader constants as cases
	switch ext {
	case ".mp3":
		return MP3Header
	case ".mp4":
		return MP3Header
	case ".wav":
		return WAVHeader
	case ".ogg":
		return OGGHeader
	case ".flac":
		return FLACHeader
	default:
		return "application/octet-stream" // unknown format
	}
}
