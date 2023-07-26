package main

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"log"
	"news-feed-bot/internal/bot"
	"news-feed-bot/internal/bot/middleware"
	"news-feed-bot/internal/botkit"
	"news-feed-bot/internal/config"
	"news-feed-bot/internal/fetcher"
	"news-feed-bot/internal/notifier"
	"news-feed-bot/internal/storage"
	"news-feed-bot/internal/summary"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot: %v", err)
		return
	}
	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}

	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier.New(
			articleStorage,
			summary.NewOpenAISummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
			botAPI,
			config.Get().NotificationInterval,
			config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.New(botAPI)
	newsBot.RegisterCommandView("start", middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViewCommandStart()))
	newsBot.RegisterCommandView("add_source", middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViewCommandAddSource(sourceStorage)))
	newsBot.RegisterCommandView("list_sources", middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViewCommandListSources(sourceStorage)))
	newsBot.RegisterCommandView("delete_source", middleware.AdminOnly(config.Get().TelegramChannelID, bot.ViewCommandDeleteSource(sourceStorage)))

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start fetcher: %v", err)
				return
			}
			log.Println("fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start notifier: %v", err)
				return
			}
			log.Println("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to run bot: %v", err)
			return
		}
		log.Println("bot stopped")
	}

}
