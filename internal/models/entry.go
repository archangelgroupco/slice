package models

// The entry model defines the stucture
// of the manifest
type Entry struct {
	MimeType      string
	RelativePath  string
	FileExtension string
	ParserVersion int8
}
