package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"messfar-telegram/domain"
	"messfar-telegram/util"
	"net/http"

	"github.com/pkg/errors"
)

type TelgramService struct {
	URL      string
	BotToken string
}

func NewTelgramService(URL, botToken string) domain.TelegramService {
	return &TelgramService{
		URL:      URL,
		BotToken: botToken,
	}
}

func (service *TelgramService) GetFileInfo(fileID string) (*domain.GetFileInfoResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/bot%s/getFile?file_id=%s", service.URL, service.BotToken, fileID))
	if err != nil {
		return nil, errors.Wrap(err, "get req fail")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body fail")
	}
	var fileInfoResponse domain.GetFileInfoResponse
	json.Unmarshal(body, &fileInfoResponse)
	return &fileInfoResponse, nil
}

func (service *TelgramService) DownloadFile(filePath string) ([]byte, error) {
	fileBytes, err := util.DownloadFile(fmt.Sprintf("%s/file/bot%s/%s", service.URL, service.BotToken, filePath))
	if err != nil {
		return nil, errors.Wrap(err, "download file failed")
	}
	return fileBytes, nil
}
