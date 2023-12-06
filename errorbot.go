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
	telegram telegram
	chatId   int
}

func New(telegramToken string, chatId int) *ErrorBot {
	return &ErrorBot{
		telegram: telegram{
			token: telegramToken,
		},
		chatId: chatId,
	}
}

func (b *ErrorBot) Send(stamp time.Time, caller string, stack string, message string) {
	go func() {
		log.Printf("Received: %s %s %s\n\n%s", stamp, caller, message, stack)
		msg := fmt.Sprintf("Message: %s\nTime: %s\nCaller: %s\nStack: %s", message, stamp, caller, stack)
		err := b.telegram.sendMessage(outgoingMessage{
			Text:   msg,
			ChatId: b.chatId,
		})
		if err != nil {
			log.Printf("Error sending: %s", err)
		}
	}()
}

func (b *ErrorBot) ZapHook(levels []zapcore.Level) func(zapcore.Entry) {
	return func(e zapcore.Entry) {
		if slices.Contains(levels, e.Level) {
			b.Send(e.Time, e.Caller.String(), e.Stack, e.Message)
		}
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
