package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/model"
)

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func ViewCommandAddSource(storage SourceStorage) botkit.ViewFunc {
	type addSourceArgs struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return err
		}
		source := model.Source{
			Name:    args.Name,
			FeedUrl: args.Url,
		}

		sourceId, err := storage.Add(ctx, source)
		if err != nil {
			return err
		}

		var (
			messageAnswer = fmt.Sprintf("Источник добавлен с ID: `%d`\\.Этот ID используется для работы с этим источником\\.", sourceId)
			reply         = tgbotapi.NewMessage(update.Message.Chat.ID, messageAnswer)
		)

		reply.ParseMode = "MarkDownV2"
		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}
