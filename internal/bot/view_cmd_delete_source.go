package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
	"strconv"
)

type SourceDeleter interface {
	Delete(ctx context.Context, id int64) error
}

func ViewCommandDeleteSource(deleter SourceDeleter) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

		args := update.Message.CommandArguments()
		sourceId, err := strconv.ParseInt(args, 10, 64)
		if err != nil {
			return err
		}

		if err := deleter.Delete(ctx, sourceId); err != nil {
			return err
		}

		var (
			messageAnswer = fmt.Sprintf("Источник с ID: `%d` удален\\.", sourceId)
			reply         = tgbotapi.NewMessage(update.Message.Chat.ID, messageAnswer)
		)

		reply.ParseMode = "MarkDownV2"
		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}
