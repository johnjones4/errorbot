package errorbot

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

type ErrorBot struct {
	application string
	telegram    telegram
	chatId      int
}

func New(application string, telegramToken string, chatId int) *ErrorBot {
	return &ErrorBot{
		application: application,
		telegram: telegram{
			token: telegramToken,
		},
		chatId: chatId,
	}
}

func (b *ErrorBot) Send(stamp time.Time, caller string, stack string, message string) {
	go func() {
		msg := fmt.Sprintf("Appliication: %s\nMessage: %s\nTime: %s\nCaller: %s\nStack: %s", b.application, message, stamp, caller, stack)
		fmt.Println(msg)
		err := b.telegram.sendMessage(outgoingMessage{
			Text:   msg,
			ChatId: b.chatId,
		})
		if err != nil {
			log.Printf("Error sending: %s", err)
		}
	}()
}

func (b *ErrorBot) ZapHook(levels []zapcore.Level) func(zapcore.Entry) error {
	return func(e zapcore.Entry) error {
		if slices.Contains(levels, e.Level) {
			b.Send(e.Time, e.Caller.String(), e.Stack, e.Message)
		}
		return nil
	}
}

type logrusHook struct {
	levels []logrus.Level
	bot    *ErrorBot
}

func (h *logrusHook) Levels() []logrus.Level {
	return h.levels
}

func (h *logrusHook) Fire(e *logrus.Entry) error {
	caller := fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)
	h.bot.Send(e.Time, caller, "", e.Message)
	return nil
}

func (b *ErrorBot) LogrusHook(levels []logrus.Level) logrus.Hook {
	return &logrusHook{
		levels: levels,
		bot:    b,
	}
}
