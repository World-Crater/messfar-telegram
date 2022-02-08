package domain

type GetFileInfoResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FileSize     int    `json:"file_size"`
		FilePath     string `json:"file_path"`
	} `json:"result"`
}

type TelegramService interface {
	GetFileInfo(fileID string) (*GetFileInfoResponse, error)
	DownloadFile(filePath string) ([]byte, error)
}
