package chatbot

import (
	"github.com/go-zoox/core-utils/fmt"

	"github.com/go-zoox/zoox/defaults"
)

// Serve starts a application server.
func (c *chatbot) Serve() error {

	app := defaults.Application()

	app.Post(c.cfg.Path, c.Handler())

	return app.Run(fmt.Sprintf(":%d", c.cfg.Port))
}
