package TelegramBotAPI

import (
	"Tasks/Service"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

type Bot struct {
	TgBot                *tgbotapi.BotAPI
	Commands             []string
	UserState            map[int64]string
	Service              *Service.Service
	BufferedSenderCard   map[int64]int
	BufferedReceiverCard map[int64]int
	personID             int
}

func InitializeBot(token string, service *Service.Service, SenderMap map[int64]int, ReceiverMap map[int64]int, UserSate map[int64]string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	return &Bot{TgBot: bot, BufferedSenderCard: SenderMap, BufferedReceiverCard: ReceiverMap, Service: service, UserState: UserSate}, nil
}
func (b *Bot) GetUpdates() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 2
	updates, err := b.TgBot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
func (b *Bot) HandlerPollingData(updatesChan tgbotapi.UpdatesChannel) {
	for update := range updatesChan {
		if update.Message == nil {
			continue
		} else if update.Message.Text == "/restart" {
			b.BufferedReceiverCard[update.Message.Chat.ID] = 0
			b.BufferedSenderCard[update.Message.Chat.ID] = 0
			b.UserState[update.Message.Chat.ID] = ""
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Cache cleaned")
			b.TgBot.Send(msg)
			continue
		} else if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, I'm your TG bot for executing transactions between wallets!")
			b.TgBot.Send(msg)
			continue
		} else if update.Message.Text == "/ping" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong!")
			b.TgBot.Send(msg)
			continue
		} else if update.Message.Text == "/transfer_money" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter your cardNumber:")
			b.UserState[update.Message.Chat.ID] = "WaitingForSenderNumber"
			b.TgBot.Send(msg)
			continue
		} else if b.UserState[update.Message.Chat.ID] == "WaitingForSenderNumber" && update.Message.Text != "" {
			number, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				b.TgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong Card Number Format, Try Again"))
				continue
			}
			c, err := b.Service.GetCardByNumber(number)
			if err != nil {
				b.TgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error fetching card:%s", err.Error())))
				continue
			}
			b.BufferedSenderCard[update.Message.Chat.ID] = number
			b.personID = c.PersonID
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter receiver cardNumber:")
			b.UserState[update.Message.Chat.ID] = "WaitingForReceiverNumber"
			b.TgBot.Send(msg)
			continue
		} else if b.UserState[update.Message.Chat.ID] == "WaitingForReceiverNumber" && update.Message.Text != "" {
			number, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				b.TgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong Card Number Format, Try Again"))
				continue
			}
			_, err = b.Service.GetCardByNumber(number)
			if err != nil {
				b.TgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Error fetching card:%s", err.Error())))
				continue
			}
			b.BufferedReceiverCard[update.Message.Chat.ID] = number
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Enter sum you want to send:")
			b.UserState[update.Message.Chat.ID] = "WaitingForSum"
			b.TgBot.Send(msg)
			continue
		} else if b.UserState[update.Message.Chat.ID] == "WaitingForSum" && update.Message.Text != "" {
			sum, err := strconv.ParseFloat(update.Message.Text, 64)
			if err != nil {
				b.TgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong sum format, Try Again"))
				continue
			}
			err = b.Service.ExecuteTransaction(b.BufferedSenderCard[update.Message.Chat.ID], b.BufferedReceiverCard[update.Message.Chat.ID], sum, context.WithValue(context.Background(), "id", b.personID))
			if err != nil {
				b.TgBot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Transaction error:%s", err.Error())))
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Transaction went successfully")
			b.TgBot.Send(msg)
			b.UserState[update.Message.Chat.ID] = ""
			b.BufferedSenderCard[update.Message.Chat.ID] = 0
			b.BufferedReceiverCard[update.Message.Chat.ID] = 0
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
		b.TgBot.Send(msg)
	}
}
