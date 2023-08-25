package models

import "time"

type Metadata struct {
	ID         uint      `json:"id"`
	Filename   string    `json:"name"`
	FileType   string    `json:"type"`
	S3Link     string    `json:"link"`
	Category   string    `json:"category"`
	UploadDate time.Time `json:"uploadDate"`
}
