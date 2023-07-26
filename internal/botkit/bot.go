package botkit

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
	"time"
)

type Bot struct {
	Api          *tgbotapi.BotAPI
	CommandViews map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *Bot {
	return &Bot{
		Api: api,
	}
}

func (b *Bot) RegisterCommandView(command string, view ViewFunc) {
	if b.CommandViews == nil {
		b.CommandViews = make(map[string]ViewFunc)
	}
	b.CommandViews[command] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.Api.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			updateCtx, updateChannel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdate(updateCtx, update)
			updateChannel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if (update.Message == nil) || (!update.Message.IsCommand()) {
		return
	}

	var view ViewFunc

	if !update.Message.IsCommand() {
		return
	}
	command := update.Message.Command()

	commandView, ok := b.CommandViews[command]
	if !ok {
		return
	}

	view = commandView
	if err := view(ctx, b.Api, update); err != nil {
		log.Printf("[ERROR] failed to handle update: %v", err)
		if _, err := b.Api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "internal error")); err != nil {
			log.Printf("[ERROR] failed to send message: %v", err)
		}
	}
}
