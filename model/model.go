package model

type FileData struct {
	ObjectName    string
	DownloadCount int64
}

type Error struct {
	Success     bool
	Description string
}
