package telegram

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/gotd/td/tg"

	"github.com/meesooqa/go-tg-bnews/internal/config"
)

// Service provides methods to interact with Telegram API
type Service struct {
	conf *config.Telegram
	api  *tg.Client
}

// NewService creates a new Telegram service with the provided API client
func NewService(conf *config.Telegram, api *tg.Client) *Service {
	return &Service{
		conf: conf,
		api:  api,
	}
}

// GetChannel retrieves a Telegram channel by its username
func (s Service) GetChannel(ctx context.Context, name string) (*tg.Channel, error) {
	resolved, err := s.api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
		Username: name,
	})
	if err != nil {
		return nil, fmt.Errorf("error resolving username %s: %w", name, err)
	}
	channel, ok := resolved.Chats[0].(*tg.Channel)
	if !ok {
		return nil, fmt.Errorf("resolved username is not a channel: %s", name)
	}
	return channel, nil
}

// GetMessages retrieves the last messages from a Telegram channel
func (s Service) GetMessages(ctx context.Context, from *tg.Channel, offsetID int) ([]*tg.Message, error) {
	messages, err := s.api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  from.ID,
			AccessHash: from.AccessHash,
		},
		Limit:    s.conf.HistoryLimit,
		OffsetID: offsetID,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting history: %w", err)
	}

	ms := messages.(*tg.MessagesChannelMessages).Messages
	result := make([]*tg.Message, 0, len(ms))
	for _, m := range ms {
		msg, ok := m.(*tg.Message)
		if !ok {
			continue
		}
		result = append(result, msg)
	}
	return result, nil
}

// ForwardMessages forwards messages from one channel to another
func (s Service) ForwardMessages(ctx context.Context, messages []*tg.Message, from, to *tg.Channel) error {
	if len(messages) == 0 {
		return nil
	}
	mID := make([]int, len(messages))
	for _, msg := range messages {
		mID = append(mID, msg.ID)
	}

	_, err := s.api.MessagesForwardMessages(ctx, &tg.MessagesForwardMessagesRequest{
		ID: mID,
		FromPeer: &tg.InputPeerChannel{
			ChannelID:  from.ID,
			AccessHash: from.AccessHash,
		},
		ToPeer: &tg.InputPeerChannel{
			ChannelID:  to.ID,
			AccessHash: to.AccessHash,
		},
		RandomID:   s.getRandomID(len(mID)),
		DropAuthor: false,
	})
	return err
}

func (Service) getRandomID(l int) []int64 {
	maxNumber := int64(math.MaxInt64)
	res := make([]int64, l)
	for i := range res {
		bi, _ := rand.Int(rand.Reader, big.NewInt(maxNumber))
		n := bi.Int64() + 1

		res[i] = n
	}
	return res
}
