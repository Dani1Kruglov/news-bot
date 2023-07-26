package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"news-feed-bot/internal/botkit"
)

func ViewCommandStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Здравствуйте, перед вами бот новостной ленты.\nДля того чтобы добавить сайт, с которого будут браться новости, введите команду /add_source, а также ссылку на сайт в таком виде:{\"name\":\"Хабр\", \"url\":\"https://habr.com/ru/rss\"}\n"+
			"Для того чтобы посмотреть список сайтов, с которых берутся новости, введите команду /list_sources\n"+
			"Для того чтобы удалить сайт, введите команду /delete_source, а также id сайтa в таком виде:{\"id\":\"2\"}")); err != nil {
			return err
		}
		return nil
	}
}
