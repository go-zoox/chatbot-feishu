package chatbot

import (
	feishuEvent "github.com/go-zoox/feishu/event"
	"github.com/go-zoox/zoox"
)

// ChatBot is the chatbot interface.
type ChatBot interface {
	Handler() zoox.HandlerFunc
	Serve() error
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
	onMessage feishuEvent.MessageHandler
}

// New creates a new chatbot
func New(cfg *Config, onMessage feishuEvent.MessageHandler) (ChatBot, error) {
	if cfg.Path == "" {
		cfg.Path = "/"
	}
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	return &chatbot{
		cfg:       cfg,
		onMessage: onMessage,
	}, nil
}
