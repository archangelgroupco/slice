package models

import "time"

// The entry model defines the stucture
// of the manifest
type Entry struct {
	MimeType      string `json:"mime_type"`
	RelativePath  string `json:"relative_path"`
	FileExtension string `json:"file_extension"`
	ParserVersion int    `json:"parser_version"`
}

type Manifest struct {
	DateTime time.Time `json:"date_time,omitempty"`
	Name     string    `json:"name,omitempty"`
	Nodes    []Entry   `json:"nodes,omitempty"`
}
