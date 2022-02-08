package domain

import "time"

type GetInfosResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Rows   []Actress
}

type GetRandomResponse []Actress

type Actress struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Romanization interface{} `json:"romanization"`
	Detail       interface{} `json:"detail"`
	Preview      string      `json:"preview"`
	Createdat    time.Time   `json:"createdat"`
	Updatedat    time.Time   `json:"updatedat"`
}

type ActressWithRecognition struct {
	Actress
	Token                 string  `json:"token"`
	RecognitionPercentage float64 `json:"recognitionPercentage"`
}

type FaceRectangle struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type PostSearchResponse []struct {
	ActressWithRecognition
}

type PostDetectResponse struct {
	RequestID string `json:"request_id"`
	TimeUsed  int    `json:"time_used"`
	Faces     []struct {
		FaceToken     string        `json:"face_token"`
		FaceRectangle FaceRectangle `json:"face_rectangle"`
	} `json:"faces"`
	ImageID string `json:"image_id"`
	FaceNum int    `json:"face_num"`
}

type PostInfosResponse struct {
	ID string `json:"id"`
}

type PostFaceResponse struct {
	FacesetToken string `json:"facesetToken"`
	FaceToken    string `json:"faceToken"`
}

type FaceService interface {
	GetInfos(limit uint, offset uint) (*GetInfosResponse, error)
	GetInfosAllActresses(limit, offsetMax int) ([]Actress, error)
	PostSearch(imageBuffer []byte) (PostSearchResponse, error)
	PostInfo(imageBuffer []byte, actress Actress) (*PostInfosResponse, error)
	PutInfo(infoID string, imageBuffer []byte) error
	PostFace(imageBuffer []byte, infoId string) (*PostFaceResponse, error)
	PostDetect(imageBuffer []byte) (*PostDetectResponse, error)
	DeleteInfo(infoID string) error
	GetRandom(quantity uint) (GetRandomResponse, error)
}
