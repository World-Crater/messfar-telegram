package util

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetMaxPhoto(photos []tgbotapi.PhotoSize) *tgbotapi.PhotoSize {
	var photo tgbotapi.PhotoSize
	for _, v := range photos {
		if photo == (tgbotapi.PhotoSize{}) || photo.FileSize < v.FileSize {
			photo = v
		}
	}
	return &photo
}
