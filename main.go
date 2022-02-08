package main

import (
	"fmt"
	"messfar-telegram/domain"
	"messfar-telegram/repo"
	"messfar-telegram/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MessageHandler struct {
	domain.TelegramService
	domain.FaceService
	Bot *tgbotapi.BotAPI
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("fatal error config. %+v\n", err))
	}
	viper.AutomaticEnv()

	telegramService := repo.NewTelgramService(viper.GetString("TELEGRAM_URL"), viper.GetString("TELEGRAM_BOT_TOKEN"))
	faceService := repo.NewFaceService(viper.GetString("FACE_SERVICE"))
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	messageHandler := MessageHandler{
		TelegramService: telegramService,
		FaceService:     faceService,
		Bot:             bot,
	}

	bot.Debug = viper.GetBool("DEBUG_MODE")

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if messageHandler.IsImageMessage(&update) {
				messageHandler.ImageMessageHandler(&update)
			} else if messageHandler.IsTextMessage(&update) {
				messageHandler.TextMessageHandler(&update)
			}
		}
	}
}

func (m *MessageHandler) IsTextMessage(update *tgbotapi.Update) bool {
	return update.Message.Text != ""
}

func (m *MessageHandler) IsImageMessage(update *tgbotapi.Update) bool {
	return len(update.Message.Photo) != 0
}

func (m *MessageHandler) TextMessageHandler(update *tgbotapi.Update) {
	switch update.Message.Text {
	case "許願":
		getRandomResponse, err := m.FaceService.GetRandom(1)
		if err != nil {
			log.Errorf("get random failed. error: %+v", err)
			return
		}

		imageBytes, err := util.DownloadFile(getRandomResponse[0].Preview)
		if err != nil {
			log.Errorf("download preview failed. error: %+v", err)
			return
		}

		textMsg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("`%s`\n點我看資料: https://messfar.com/?ID=%s", getRandomResponse[0].Name, getRandomResponse[0].ID),
		)
		textMsg.ReplyToMessageID = update.Message.MessageID
		photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
			Name:  "image",
			Bytes: imageBytes,
		})
		photoMsg.ReplyToMessageID = update.Message.MessageID

		if _, err := m.Bot.Send(textMsg); err != nil {
			log.Errorf("send text message failed. error: %+v", err)
		}
		if _, err := m.Bot.Send(photoMsg); err != nil {
			log.Errorf("send photo message failed. error: %+v", err)
		}
	case "我心愛的女孩":
		textMsg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"請至髒沙發頁面觀看: https://messfar.com",
		)
		if _, err := m.Bot.Send(textMsg); err != nil {
			log.Errorf("send text message failed. error: %+v", err)
		}
	}
}

func (m *MessageHandler) ImageMessageHandler(update *tgbotapi.Update) {
	getFileInfoResponse, err := m.TelegramService.GetFileInfo(util.GetMaxPhoto(update.Message.Photo).FileID)
	if err != nil {
		log.Errorf("get file info failed. error: %+v", err)
		return
	}

	downloadFileBytes, err := m.TelegramService.DownloadFile(getFileInfoResponse.Result.FilePath)
	if err != nil {
		log.Errorf("download file info failed. error: %+v", err)
		return
	}

	postSearchResponse, err := m.FaceService.PostSearch(downloadFileBytes)
	if err != nil {
		log.Errorf("search face failed. error: %+v", err)
		return
	}

	for _, v := range postSearchResponse {
		imageBytes, err := util.DownloadFile(v.Actress.Preview)
		if err != nil {
			log.Errorf("download preview failed. error: %+v", err)
			return
		}

		textMsg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf(
				"我猜可能是`%s`\n相似度: %f%%\n點我看資料: https://messfar.com/?ID=%s",
				v.Actress.Name,
				v.RecognitionPercentage,
				v.Actress.ID,
			),
		)
		textMsg.ReplyToMessageID = update.Message.MessageID
		photoMsg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
			Name:  "image",
			Bytes: imageBytes,
		})
		photoMsg.ReplyToMessageID = update.Message.MessageID

		if _, err := m.Bot.Send(textMsg); err != nil {
			log.Errorf("send text message failed. error: %+v", err)
		}
		if _, err := m.Bot.Send(photoMsg); err != nil {
			log.Errorf("send photo message failed. error: %+v", err)
		}
	}
}
