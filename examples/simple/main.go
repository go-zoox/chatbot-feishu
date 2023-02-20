package simple

import (
	"github.com/go-zoox/chatbot-feishu"
	"github.com/go-zoox/feishu/event"
	"github.com/go-zoox/logger"
)

func main() {
	bot, err := chatbot.New(&chatbot.Config{
		AppID:     "xxx",
		AppSecret: "yyy",
	})
	if err != nil {
		logger.Errorf("failed to create bot: %v", err)
		return
	}

	bot.OnCommand("/chatgpt", &chatbot.Command{
		ArgsLength: 2,
		Handler: func(args []string, request *event.EventRequest, reply func(content string, msgType ...string) error) error {
			return nil
		},
	})

	if err := bot.Run(); err != nil {
		logger.Fatalf("%v", err)
	}
}
