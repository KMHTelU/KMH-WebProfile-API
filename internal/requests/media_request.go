package requests

type CreateMediaRequest struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
	AltText  string `json:"alt_text"`
	Caption  string `json:"caption"`
}
