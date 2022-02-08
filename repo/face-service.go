package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"messfar-telegram/domain"
	"mime/multipart"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type FaceService struct {
	Url string
}

func NewFaceService(Url string) domain.FaceService {
	return &FaceService{
		Url: Url,
	}
}

func (service *FaceService) GetRandom(quantity uint) (domain.GetRandomResponse, error) {
	if quantity == 0 {
		quantity = 1
	}

	res, err := http.Get(fmt.Sprintf("%s/faces/random?quantity=%d", service.Url, quantity))
	if err != nil {
		return nil, errors.Wrap(err, "get req fail")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body fail")
	}
	var getRandomResponse domain.GetRandomResponse
	json.Unmarshal(body, &getRandomResponse)
	return getRandomResponse, nil
}

func (service *FaceService) GetInfos(limit uint, offset uint) (*domain.GetInfosResponse, error) {
	if limit == 0 {
		return nil, errors.New("require limit")
	}

	res, err := http.Get(fmt.Sprintf("%s/faces/infos?limit=%d&offset=%d", service.Url, limit, offset))
	if err != nil {
		return nil, errors.Wrap(err, "get req fail")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body fail")
	}
	var infosResponse domain.GetInfosResponse
	json.Unmarshal(body, &infosResponse)
	return &infosResponse, nil
}

func (service *FaceService) GetInfosAllActresses(limit, count int) ([]domain.Actress, error) {
	actresses := []domain.Actress{}
	offset := 0
	if limit == -1 {
		limit = 1000
	}

	for {
		GetInfosResponse, err := service.GetInfos(uint(limit), uint(offset))
		if err != nil {
			return nil, errors.Wrap(err, "get infos fail")
		}
		if count == -1 {
			count = GetInfosResponse.Count
		}
		actresses = append(actresses, GetInfosResponse.Rows...)
		offset = offset + limit
		if offset >= count {
			log.Info("current actress infos count: ", len(actresses))
			return actresses, nil
		}
	}
}

func (service *FaceService) createImagePayload(imageBuffer []byte, keyName string) (*bytes.Buffer, *multipart.Writer, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	imagePart, err := writer.CreateFormFile(keyName, keyName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "create form file failed")
	}

	if _, err = imagePart.Write(imageBuffer); err != nil {
		return nil, nil, errors.Wrap(err, "io copy failed")
	}
	return payload, writer, nil
}

func (service *FaceService) PostDetect(imageBuffer []byte) (*domain.PostDetectResponse, error) {
	payload, writer, err := service.createImagePayload(imageBuffer, "image")
	if err != nil {
		return nil, errors.Wrap(err, "create image payload failed")
	}
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "close writer failed")
	}
	res, err := http.Post(fmt.Sprintf("%s/faces/detect", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "post detect API failed")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal server error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body failed")
	}
	var postDetectResponse domain.PostDetectResponse
	json.Unmarshal(body, &postDetectResponse)
	return &postDetectResponse, nil
}

func (service *FaceService) PostSearch(imageBuffer []byte) (domain.PostSearchResponse, error) {
	payload, writer, err := service.createImagePayload(imageBuffer, "image")
	if err != nil {
		return nil, errors.Wrap(err, "create image payload failed")
	}
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "close writer failed")
	}
	res, err := http.Post(fmt.Sprintf("%s/faces/search", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "post search API failed")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal server error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body failed")
	}
	var postSearchResponse domain.PostSearchResponse
	json.Unmarshal(body, &postSearchResponse)
	return postSearchResponse, nil
}

func (service *FaceService) DeleteInfo(infoID string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/faces/info/%s", service.Url, infoID), nil)
	if err != nil {
		return errors.Wrap(err, "create delete request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "delete info failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("get %d error", res.StatusCode))
	}

	return nil
}

func (service *FaceService) PostInfo(imageBuffer []byte, actress domain.Actress) (*domain.PostInfosResponse, error) {
	if actress.Name == "" {
		return nil, errors.New("actress name is empty")
	}

	payload, writer, err := service.createImagePayload(imageBuffer, "preview")
	if err != nil {
		return nil, errors.Wrap(err, "create image payload failed")
	}
	_ = writer.WriteField("name", actress.Name)
	if actress.Romanization != "" && actress.Romanization != nil {
		if err := writer.WriteField("romanization", actress.Romanization.(string)); err != nil {
			return nil, errors.Wrap(err, "write field 'romanization' failed")
		}
	}
	if actress.Detail != "" && actress.Detail != nil {
		if err := writer.WriteField("detail", actress.Detail.(string)); err != nil {
			return nil, errors.Wrap(err, "write field 'detail' failed")
		}
	}
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "close writer failed")
	}

	res, err := http.Post(fmt.Sprintf("%s/faces/info", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "post info failed")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body failed")
	}
	var postInfosResponse domain.PostInfosResponse
	json.Unmarshal(body, &postInfosResponse)
	return &postInfosResponse, nil
}

func (service *FaceService) PutInfo(infoID string, imageBuffer []byte) error {
	payload, writer, err := service.createImagePayload(imageBuffer, "preview")
	if err != nil {
		return errors.Wrap(err, "create image payload failed")
	}
	if err := writer.Close(); err != nil {
		return errors.Wrap(err, "close writer failed")
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/faces/info/%s", service.Url, infoID), payload)
	if err != nil {
		return errors.Wrap(err, "create delete request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "delete info failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("put info error. status %d", res.StatusCode))
	}

	return nil
}

func (service *FaceService) PostFace(imageBuffer []byte, infoId string) (*domain.PostFaceResponse, error) {
	if infoId == "" {
		return nil, errors.New("require infoId")
	}

	payload, writer, err := service.createImagePayload(imageBuffer, "image")
	if err != nil {
		return nil, errors.Wrap(err, "create image payload failed")
	}

	_ = writer.WriteField("infoId", infoId)

	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "close writer failed")
	}

	res, err := http.Post(fmt.Sprintf("%s/faces/face", service.Url), writer.FormDataContentType(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "post face API failed")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, errors.New("internal error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body failed")
	}
	var postFaceResponse domain.PostFaceResponse
	json.Unmarshal(body, &postFaceResponse)
	return &postFaceResponse, nil
}
