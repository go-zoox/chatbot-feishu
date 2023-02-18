package chatbot

import (
	"net/http"

	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/feishu"
	feishuEvent "github.com/go-zoox/feishu/event"
	"github.com/go-zoox/zoox"
)

// Handler creates a chat bot handler for zoox.
func (c *chatbot) Handler() zoox.HandlerFunc {
	bot := feishu.New(&feishu.Config{
		AppID:     c.cfg.AppID,
		AppSecret: c.cfg.AppSecret,
	})

	return func(ctx *zoox.Context) {
		var request feishuEvent.EventRequest
		if err := ctx.BindJSON(&request); err != nil {
			ctx.Fail(err, 500, "Internal Server Error")
			return
		}

		if request.IsChallenge() {
			ctx.Logger.Infof("challenge request => %s", request.Challenge)

			if request.Challenge == "" {
				ctx.Fail(fmt.Errorf("expect challenge, but got empty"), 400000, "expect challenge, but got empty")
				return
			}

			ctx.JSON(http.StatusOK, zoox.H{
				"challenge": request.Challenge,
			})

			return
		}

		event := bot.Event(&request)

		// go event.OnChatReceiveMessage(func(contentString string, request *feishuEvent.EventRequest, reply func(content string) error) error {
		// 	if contentString != "" {
		// 		type Content struct {
		// 			Text string `json:"text"`
		// 		}
		// 		var content Content
		// 		if err := json.Unmarshal([]byte(contentString), &content); err != nil {
		// 			return err
		// 		}
		// 	}

		// 	return nil
		// })

		go event.OnChatReceiveMessage(c.onMessage)

		ctx.Success(nil)
	}
}
