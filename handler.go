package chatbot

import (
	"encoding/json"
	"net/http"

	"github.com/go-zoox/core-utils/strings"
	"github.com/go-zoox/debug"
	"github.com/go-zoox/logger"

	"github.com/go-zoox/chatbot-feishu/command"
	"github.com/go-zoox/chatbot-feishu/utils"
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

		if request.Encrypt != "" {
			if c.cfg.EncryptKey == "" {
				ctx.Fail(fmt.Errorf("encrypt key missing"), 400000, "encrypt key missing")
				return
			}

			plaintext, err := utils.Decrypt(request.Encrypt, c.cfg.EncryptKey)
			if err != nil {
				ctx.Fail(fmt.Errorf("failed to descript, please check encrypt key: %v", err), 400001, "failed to descript, please check encrypt key")
				return
			}

			if err := json.Unmarshal(plaintext, &request); err != nil {
				ctx.Fail(fmt.Errorf("failed to parse descripted body: %v", err), 400002, "failed to parse descripted body")
				return
			}
		}

		if request.IsChallenge() {
			ctx.Logger.Infof("challenge request => %s", request.Challenge)

			if c.cfg.VerificationToken != "" {
				if request.Token != c.cfg.VerificationToken {
					logger.Infof("verification token expect %s, but got %s", c.cfg.VerificationToken, request.Token)
					ctx.Fail(fmt.Errorf("verification tokens are not matched"), 400001, "verification tokens are not matched")
					return
				}
			}

			ctx.JSON(http.StatusOK, zoox.H{
				"challenge": request.Challenge,
			})

			return
		}

		event := bot.Event(&request)

		go func() {
			err := event.OnChatReceiveMessage(func(content string, request *feishuEvent.EventRequest, reply MessageReply) error {
				if debug.IsDebugMode() {
					fmt.PrintJSON(request)
				}

				if request.Event.Message.MessageType != "text" {
					logger.Infof("ignore message type: %s", request.Event.Message.MessageType)
					return nil
				}

				type Content struct {
					Text string `json:"text"`
				}
				var contentX Content
				if err := json.Unmarshal([]byte(content), &contentX); err != nil {
					return fmt.Errorf("failed to parse message content(%s): %v", content, err)
				}

				if contentX.Text == "" {
					logger.Infof("ignore empty message: %s", content)
					return nil
				}

				logger.Infof("message: %s", contentX.Text)

				// if len(c.events) != 0 {
				// 	if event, ok := c.events[request.EventType()]; ok {
				// 		if err := event.Handler(request, reply); err != nil {
				// 			logger.Warn("failed to listen event(%s): %v", request.EventType(), err)
				// 		}
				// 	}
				// }

				if len(c.commands) != 0 {
					text := contentX.Text
					if strings.StartsWith(text, "@_user_1 ") {
						text = text[len("@_user_1 "):]
					}
					text = strings.TrimSpace(text)

					logger.Infof("start to check whether %s is a command ...", contentX.Text)
					if command.IsCommand(text) {
						cmd, arg, err := command.ParseCommandWithArg(text)
						if err != nil {
							return fmt.Errorf("failed to parse command(%s): %v", text, err)
						}

						logger.Infof("start to check whether command(cmd: %s, arg: %s) exists ...", cmd, arg)
						if c, ok := c.commands[cmd]; ok {
							var args []string
							if c.ArgsLength != 0 {
								args = strings.SplitN(arg, " ", c.ArgsLength)
							} else {
								args = []string{arg}
							}

							if err := c.Handler(args, request, reply); err != nil {
								logger.Errorf("failed to run command(%s): %v", contentX.Text, err)
								reply("failed to run command %s: %s", cmd, err.Error())
							}
							return nil
						}
					}
				}

				logger.Infof("fallback to common message: %s...", content)
				if c.onMessage != nil {
					if err := c.onMessage(contentX.Text, request, reply); err != nil {
						logger.Warn("failed to list message: %v", err)
					}
				}

				return nil
			})
			if err != nil {
				logger.Errorf("failed to OnChatReceiveMessage: %v", err)
			}
		}()

		ctx.Success(nil)
	}
}
