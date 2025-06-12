package proc

import (
	"fmt"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"

	mytg "github.com/meesooqa/go-tg-bnews/internal/telegram"
)

// AuthProcessor authenticates the client using the provided auth flow
func AuthProcessor(flow auth.Flow) Processor {
	return func(st *PipelineState) error {
		if err := st.Client.Auth().IfNecessary(st.Ctx, flow); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
		return nil
	}
}

// InitServiceProcessor initializes the service with the Telegram API client
func InitServiceProcessor() Processor {
	return func(st *PipelineState) error {
		st.Service = mytg.NewService(st.Conf.Telegram, st.Client.API())
		st.chanCache = make(map[string]*tg.Channel)
		return nil
	}
}

// LoadChannelsProcessor loads channels by their names and caches them
func LoadChannelsProcessor(names ...string) Processor {
	return func(st *PipelineState) error {
		st.cacheMu.Lock()
		defer st.cacheMu.Unlock()

		for _, name := range names {
			if _, ok := st.chanCache[name]; !ok {
				ch, err := st.Service.GetChannel(st.Ctx, name)
				if err != nil {
					return fmt.Errorf("get channel %s: %w", name, err)
				}
				st.chanCache[name] = ch
			}
		}
		st.ChannelFrom = st.chanCache[names[0]]
		st.ChannelTo = st.chanCache[names[1]]
		return nil
	}
}

// FetchMessagesProcessor fetches messages from the source channel
func FetchMessagesProcessor() Processor {
	return func(st *PipelineState) error {
		msgs, err := st.Service.GetMessages(st.Ctx, st.ChannelFrom)
		if err != nil {
			return fmt.Errorf("fetch messages: %w", err)
		}
		st.Messages = msgs
		if len(msgs) == 0 {
			return fmt.Errorf("no messages in %s", st.ChannelFrom.Username)
		}
		return nil
	}
}

// FilterProcessor applies the provided filters to the messages in the pipeline state
func FilterProcessor(filters ...MessageFilter) Processor {
	return func(st *PipelineState) error {
		st.Messages = FilterMessages(st.Messages, filters...)
		return nil
	}
}

// ForwardProcessor forwards the messages from the source channel to the destination channel
func ForwardProcessor() Processor {
	return func(st *PipelineState) error {
		if len(st.Messages) == 0 {
			return nil
		}
		return st.Service.ForwardMessages(st.Ctx, st.Messages, st.ChannelFrom, st.ChannelTo)
	}
}
