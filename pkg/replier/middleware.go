package replier

import (
	"context"

	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
)

var replierKey = replierKeyType{}

type replierKeyType struct{}

// NewMiddleware creates a new middleware for the passed *Replier, that stores
// the replier in the *plugin.Context.
func NewMiddleware(o *Options) bot.MiddlewareFunc {
	o.fillDefaults()

	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, pctx *plugin.Context) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			r := &replier{
				Options: o,
				s:       s,
				ctx:     ctx,
				cancel:  cancel,
			}

			pctx.Set(replierKey, r)
			return next(s, pctx)
		}
	}
}

func get(ctx *plugin.Context) *replier { return ctx.Get(replierKey).(*replier) }
