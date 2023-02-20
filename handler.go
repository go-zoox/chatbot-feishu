package chatbot

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/logger"

	"github.com/go-zoox/chatbot-feishu/command"
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

		go event.OnChatReceiveMessage(func(content string, request *feishuEvent.EventRequest, reply func(content string, msgType ...string) error) error {
			if c.onMessage != nil {
				if err := c.onMessage(content, request, reply); err != nil {
					logger.Warn("failed to list message: %v", err)
				}
			}

			if len(c.events) != 0 {
				if event, ok := c.events[request.EventType()]; ok {
					if err := event.Handler(request, reply); err != nil {
						logger.Warn("failed to listen event(%s): %v", request.EventType(), err)
					}
				}
			}

			if len(c.commands) != 0 {
				type Content struct {
					Text string `json:"text"`
				}

				var contentX Content
				if err := json.Unmarshal([]byte(content), &content); err != nil {
					return err
				}

				if !command.IsCommand(contentX.Text) {
					return nil
				}

				cmd, arg, err := command.ParseCommandWithArg(contentX.Text)
				if err != nil {
					return fmt.Errorf("failed to parse command(%s): %v", contentX.Text, err)
				}

				if c, ok := c.commands[cmd]; ok {
					if err := c.Handler(strings.SplitN(arg, " ", c.ArgsLength), request, reply); err != nil {
						logger.Errorf("failed to run command(%s): %v", contentX.Text, err)
						reply("failed to run command %s: %s", cmd, err.Error())
					}
				}
			}

			return nil
		})

		ctx.Success(nil)
	}
}
