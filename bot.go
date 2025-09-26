package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	ChatID int64 `json:"chat_id"`
	IsSub  bool  `json:"is_sub"`
}

var userlist = []User{}
var dict = make(map[int64]bool)

func main() {
	Initialize()

	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		fmt.Println(err)
	}

	commands := []tgbotapi.BotCommand{
		{Command: "/sub", Description: "подписаться на отправку и получение тиктоков другим пользователям"},
		{Command: "/unsub", Description: "отписаться от получения и отправки тиктоков другим пользователям"},
	}
	setCommandsConfig := tgbotapi.NewSetMyCommands(commands...)
	if _, err := bot.Request(setCommandsConfig); err != nil {
		fmt.Print("Ошибка при установке команд меню: ", err)
	} else {
		fmt.Print("Команды меню успешно установлены!")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		us := User{ChatID: update.Message.Chat.ID, IsSub: false}
		isUserHere(us)

		command := update.Message.Command()

		switch command {
		case "sub":
			Unmarshall()
			var chatID int64
			chatID = update.Message.Chat.ID
			for i := 0; i < len(userlist); i++ {
				if chatID == userlist[i].ChatID {
					if !userlist[i].IsSub {
						dict[chatID] = true
						userlist[i].IsSub = true
					}
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно подписан на отправку и получение тикитоков")
			bot.Send(msg)
			Marshall()
		case "unsub":
			Unmarshall()
			var chatID int64
			chatID = update.Message.Chat.ID
			for i := 0; i < len(userlist); i++ {
				if chatID == userlist[i].ChatID {
					if userlist[i].IsSub {
						dict[chatID] = false
						userlist[i].IsSub = false
					}
				}
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Успешно отписан от отправки и получения тикитоков")
			bot.Send(msg)
			Marshall()
		}

		if update.Message != nil && !update.Message.IsCommand() {
			err := Download(update.Message.Text, "downloads", "C:/Users/Qwerty/scoop/apps/yt-dlp/2025.04.30/yt-dlp.exe")
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				bot.Send(msg)
			}
			fp := "downloads/1.mp4"
			msg := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FilePath(fp))
			bot.Send(msg)
			var userChatID = update.Message.Chat.ID
			for chatID := range dict {
				if chatID == userChatID {
					continue
				}
				if dict[chatID] {
					msg := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(fp))
					bot.Send(msg)
				}
			}
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

func Marshall() {
	bytes, err := json.Marshal(userlist)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.Create("users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, er := file.Write(bytes)
	if er != nil {
		fmt.Println(er)
	}
}

func Unmarshall() {
	var us []User
	data, err := os.ReadFile("users.json")
	if err != nil {
		fmt.Println(err)
	}
	if len(data) == 0 {
		return
	} else {
		er := json.Unmarshal(data, &us)
		if er != nil {
			fmt.Println(er)
		}
		userlist = us
	}
}

func DictWrite() {
	Unmarshall()
	for i := 0; i < len(userlist); i++ {
		dict[userlist[i].ChatID] = userlist[i].IsSub
	}
	Marshall()
}

func Initialize() {
	Unmarshall()
	Marshall()
	DictWrite()
}

func isUserHere(us User) {
	Unmarshall()
	isHere := false
	for i := 0; i < len(userlist); i++ {
		if us.ChatID == userlist[i].ChatID {
			isHere = true
			break
		}
	}
	if !isHere {
		userlist = append(userlist, us)
	}
	Marshall()
}
