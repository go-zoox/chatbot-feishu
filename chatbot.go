package chatbot

import (
	"github.com/go-zoox/core-utils/fmt"
	feishuEvent "github.com/go-zoox/feishu/event"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

// OnMessageHandler ...
type OnMessageHandler = feishuEvent.MessageHandler

// MessageReply ...
type MessageReply = func(content string, msgType ...string) error

// Command ...
type Command struct {
	ArgsLength int `json:"args_length,omitempty"`
	Handler    func(args []string, request *feishuEvent.EventRequest, reply MessageReply) error
}

// Event ...
type Event struct {
	Handler func(request *feishuEvent.EventRequest, reply MessageReply) error
}

// ChatBot is the chatbot interface.
type ChatBot interface {
	OnEvent(event string, handler *Event) error
	OnMessage(handler OnMessageHandler) error
	OnCommand(command string, handler *Command) error
	Run() error
	//
	Handler() zoox.HandlerFunc
}

// Config is the configuration for create chatbot.
type Config struct {
	ChatGPTAPIKey     string
	AppID             string
	AppSecret         string
	EncryptKey        string
	VerificationToken string
	//
	Port int64
	Path string
}

type chatbot struct {
	cfg       *Config
	onMessage OnMessageHandler
	events    map[string]*Event
	commands  map[string]*Command
}

// New creates a new chatbot
func New(cfg *Config) (ChatBot, error) {
	if cfg.Path == "" {
		cfg.Path = "/"
	}
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	return &chatbot{
		cfg:      cfg,
		events:   make(map[string]*Event),
		commands: map[string]*Command{},
	}, nil
}

func (c *chatbot) OnMessage(handler OnMessageHandler) error {
	if c.onMessage != nil {
		return fmt.Errorf("on message is already registered")
	}

	c.onMessage = handler
	return nil
}

func (c *chatbot) OnEvent(event string, handler *Event) error {
	if _, ok := c.events[event]; ok {
		return fmt.Errorf("failed to register event %s, which is already registered before", event)
	}

	c.events[event] = handler
	return nil
}

func (c *chatbot) OnCommand(command string, handler *Command) error {
	if _, ok := c.commands[command]; ok {
		return fmt.Errorf("failed to register command %s, which is already registered before", command)
	}

	logger.Infof("register command: %s", command)
	c.commands[command] = handler
	return nil
}

// Run starts a application server.
func (c *chatbot) Run() error {
	app := defaults.Application()

	app.Post(c.cfg.Path, c.Handler())

	return app.Run(fmt.Sprintf(":%d", c.cfg.Port))
}
