// Package replier provides the delayed replier.
package replier

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
)

type replier struct {
	*Options

	s      *state.State
	ctx    context.Context
	cancel func()
}

// Cancel cancels the replier, making ongoing calls to one of the Reply methods
// return errors.Abort.
func Cancel(ctx *plugin.Context) { get(ctx).cancel() }

// Reply sends a textual reply in the invoking channel.
// The passed content may be split in accordance with the SplitterFunc defined
// in the replier's Options.
func Reply(ctx *plugin.Context, content string) ([]discord.Message, error) {
	return ReplyWithDelay(ctx, content, nil)
}

// ReplyWithDelay is the same as Reply, but uses the passed DelayFunc
// instead of the replier's default.
func ReplyWithDelay(ctx *plugin.Context, content string, f DelayFunc) ([]discord.Message, error) {
	r := get(ctx)

	if f == nil {
		f = r.DefaultDelayFunc
	}

	return r.reply(r.s, ctx, ctx.ChannelID, ctx.ReplyMessage, content, f)
}

// ReplyDM sends a textual reply in a direct message channel with the invoking
// user.
// The passed content may be split in accordance with the SplitterFunc defined
// in the replier's Options.
func ReplyDM(ctx *plugin.Context, content string) ([]discord.Message, error) {
	return ReplyDMWithDelay(ctx, content, nil)
}

// ReplyDMWithDelay is the same as Reply, but uses the passed DelayFunc
// instead of the replier's default.
func ReplyDMWithDelay(ctx *plugin.Context, content string, f DelayFunc) ([]discord.Message, error) {
	r := get(ctx)

	dm, err := r.s.CreatePrivateChannel(ctx.User.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if f == nil {
		f = r.DefaultDelayFunc
	}

	return r.reply(r.s, ctx, dm.ID, ctx.ReplyMessageDM, content, f)
}

type sendFunc func(data api.SendMessageData) (*discord.Message, error)

func (r *replier) reply(
	s *state.State, pctx *plugin.Context, channelID discord.ChannelID, f sendFunc, content string,
	delayFunc DelayFunc,
) ([]discord.Message, error) {
	select {
	case <-r.ctx.Done():
		return nil, errors.Abort
	default:
	}

	contents := r.Splitter(content)

	msgs := make([]discord.Message, len(contents))

	for _, c := range contents {
		stopTyping := func() {}
		if !r.NoTyping {
			var ctx context.Context

			ctx, stopTyping = context.WithCancel(r.ctx)
			go startTyping(ctx, s, pctx, channelID)
		}

		select {
		case <-r.ctx.Done():
			stopTyping()
			return msgs, errors.Abort
		case <-time.After(delayFunc(c)):
		}

		msg, err := f(api.SendMessageData{Content: c})
		stopTyping()
		if err != nil {
			return msgs, nil
		}

		msgs = append(msgs, *msg)
	}

	return msgs, nil
}
