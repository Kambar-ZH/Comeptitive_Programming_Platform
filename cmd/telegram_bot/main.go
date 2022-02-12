package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"site/internal/config"
	"site/internal/dto"
	"site/internal/logger"
	"site/internal/store/postgres"
	"time"

	telebot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func roundTimeUp(t time.Time) time.Time {
	rounded := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	return rounded
}

func main() {
	bot, err := telebot.NewBotAPI(config.TelebotToken())
	if err != nil {
		log.Panic(err)
	}
	db := postgres.NewDB()
	if err := db.Connect(config.DSN()); err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
	endOfDay := roundTimeUp(time.Now())
	logger.Logger.Sugar().Debugf("%d %d", time.Now(), endOfDay)

	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
	for {
		time.Sleep(5 * time.Second)
		contests, err := db.Contests().FindByTimeInterval(context.Background(), &dto.ContestFindByTimeInterval{
			StartTime: time.Now(),
			EndTime:   endOfDay,
		})
		if err != nil {
			log.Panic(err)
		}
		for _, contest := range contests {
			if time.Until(contest.StartDate.Add(-1*time.Hour)) < 0 {
				msg := telebot.NewMessageToChannel(
					config.TelebotChannelName(),
					fmt.Sprintf("%s starts in ~ 1 hour", contest.Name),
				)
				str, err := bot.Send(msg)
				if err != nil {
					logger.Logger.Error(err.Error())
				}
				logger.Logger.Sugar().Debugf("%v", str)
			}
		}
	}
}