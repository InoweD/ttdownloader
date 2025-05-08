package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var subs []int64 = []int64{5825066447}

func main() {
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		fmt.Println(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message != nil {
			err := Download(update.Message.Text, "downloads", "C:/Users/Qwerty/scoop/apps/yt-dlp/2025.04.30/yt-dlp.exe")
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				bot.Send(msg)
			}
			fp := "downloads/1.mp4"
			msg := tgbotapi.NewVideo(subs[0], tgbotapi.FilePath(fp))
			bot.Send(msg)
			os.Remove(fp)
		}
	}
}

func Download(tiktokURL string, outputDir string, ytdlpPath string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Printf("Создание директории: %s", outputDir)
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Println(err)
		}
	}

	if ytdlpPath == "" {
		ytdlpPath = "yt-dlp"
		log.Println("Путь к yt-dlp не указан, используется 'yt-dlp' из PATH.")
	} else {
		log.Printf("Используется указанный путь к yt-dlp: %s", ytdlpPath)
	}

	outputPathTemplate := filepath.Join(outputDir, "1.%(ext)s")

	cmd := exec.Command(ytdlpPath, "-o", outputPathTemplate, tiktokURL)
	outputBytes, err := cmd.CombinedOutput()
	outputString := string(outputBytes)
	if err != nil {
		return err
	}

	log.Printf("Вывод yt-dlp:\n%s", outputString)
	log.Printf("Видео успешно скачано: %s", tiktokURL)
	return nil
}
