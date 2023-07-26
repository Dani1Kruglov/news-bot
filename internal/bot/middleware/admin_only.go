package middleware

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
)

func AdminOnly(channelId int64, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelId,
				},
			})

		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)

			}
		}
		if _, err := bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"К сожалению, у вас нет доступа к таким командам.",
		)); err != nil {
			return err
		}
		return nil
	}
}
