package util

import (
	"bytes"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
)

var (
	isImageReg *regexp.Regexp
)

func Init() error {
	reg, err := regexp.Compile(`(\.png|\.jpeg|\.jpg)$`)
	if err != nil {
		return errors.Wrap(err, "create regexp failed")
	}
	isImageReg = reg
	return nil
}

func IsImage(url string) bool {
	return isImageReg.MatchString(url)
}

func DownloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request fail")
	}
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request fail")
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(response.Body); err != nil {
		return nil, errors.Wrap(err, "decode file failed")
	}

	return buf.Bytes(), nil
}
