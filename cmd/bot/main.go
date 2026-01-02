package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"time"

	"telegram-captcha-bot/internal/captcha"
	"telegram-captcha-bot/internal/storage"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TG_TOKEN")
	if token == "" {
		log.Fatal("TG_TOKEN is not set")
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.NewStorage(3 * time.Minute)

	b.Handle("/start", func(c tele.Context) error {
		code, imgData, err := captcha.Generate()
		if err != nil {
			return c.Send("Internal error generating captcha.")
		}

		store.Set(c.Sender().ID, code)

		photo := &tele.Photo{
			File:    tele.FromReader(bytes.NewReader(imgData)),
			Caption: "Введите символы с картинки:",
		}
		return c.Send(photo)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		text := strings.ToUpper(strings.TrimSpace(c.Text()))
		if text == "" || text == "/START" {
			return nil
		}

		if store.Verify(c.Sender().ID, text) {
			return c.Send("Верно! Вы успешно прошли проверку.")
		}

		return c.Send("Неверно или время истекло. Введите /start для новой попытки.")
	})

	log.Println("Bot started")
	b.Start()
}
