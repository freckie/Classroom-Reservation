package models

type FilesGetResponse struct {
	FilesCount int            `json:"files_count"`
	Files      []FilesGetItem `json:"files"`
}

type FilesGetItem struct {
	FileID    string `json:"file_id"`
	FileName  string `json:"file_name"`
	CreatedAt string `json:"created_at"`
}

type FilesPostResponse struct {
	FileID string `json:"file_id"`
}
