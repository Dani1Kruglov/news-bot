package notifier

import (
	"context"
	"fmt"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"news-feed-bot/internal/botkit/markup"
	"news-feed-bot/internal/model"
	"regexp"
	"strings"
	"time"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkPosted(ctx context.Context, id int64) error
}

type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

type Notifier struct {
	Articles         ArticleProvider
	Summarizer       Summarizer
	Bot              *tgbotapi.BotAPI
	SendInterval     time.Duration
	LookupTimeWindow time.Duration
	ChannelId        int64
}

func New(
	articleProvider ArticleProvider,
	summarizer Summarizer,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelID int64,
) *Notifier {
	return &Notifier{
		Articles:         articleProvider,
		Summarizer:       summarizer,
		Bot:              bot,
		SendInterval:     sendInterval,
		LookupTimeWindow: lookupTimeWindow,
		ChannelId:        channelID,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.SendInterval)
	defer ticker.Stop()
	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticles, err := n.Articles.AllNotPosted(ctx, time.Now().Add(-n.LookupTimeWindow), 1)
	if err != nil {
		return err
	}
	if len(topOneArticles) == 0 {
		return nil
	}
	article := topOneArticles[0]
	summary, err := n.extractSummary(ctx, article)
	if err != nil {
		return err
	}
	if err := n.sendArticle(article, summary); err != nil {
		return err
	}

	return n.Articles.MarkPosted(ctx, article.Id)

}

func (n *Notifier) extractSummary(ctx context.Context, article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}
	summary, err := n.Summarizer.Summarize(ctx, cleanText(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

func cleanText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(article model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"
	msg := tgbotapi.NewMessage(n.ChannelId, fmt.Sprintf(msgFormat, markup.EscapeForMarkdown(article.Title), markup.EscapeForMarkdown(summary), markup.EscapeForMarkdown(article.Link)))
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := n.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
