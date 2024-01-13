package file

import "time"

type FileInformation struct {
	ID         uint      `json:"id"`
	Filename   string    `json:"name"`
	FileType   string    `json:"type"`
	S3Link     string    `json:"link"`
	Category   string    `json:"category"`
	UploadDate time.Time `json:"uploadDate"`
	Metadata   `json:"metadata"`
}

type Metadata struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Year   int `json:"year"`
}

