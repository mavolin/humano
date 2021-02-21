package replier

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
)

func startTyping(ctx context.Context, s *state.State, pctx *plugin.Context, channelID discord.ChannelID) {
	t := time.NewTicker(6 * time.Second)

	err := s.Typing(channelID)
	if err != nil {
		pctx.HandleErrorSilently(err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			_ = s.Typing(channelID)
			if err != nil {
				pctx.HandleErrorSilently(err)
				return
			}
		}
	}
}
