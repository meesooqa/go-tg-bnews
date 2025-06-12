package proc

import (
	"context"
	"sync"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/meesooqa/go-tg-bnews/internal/config"
	mytg "github.com/meesooqa/go-tg-bnews/internal/telegram"
)

// PipelineState holds the state of the processing pipeline
type PipelineState struct {
	Ctx    context.Context
	Client *telegram.Client
	Conf   *config.AppConfig

	Service          *mytg.Service
	ChannelFrom      *tg.Channel
	ChannelTo        *tg.Channel
	Messages         []*tg.Message
	messagesOffsetID int

	chanCache map[string]*tg.Channel
	cacheMu   sync.Mutex
}

// Processor defines a function type that processes the PipelineState
type Processor func(*PipelineState) error

// Chain creates a new Processor that chains multiple processors together
func Chain(stages ...Processor) Processor {
	return func(st *PipelineState) error {
		for _, stage := range stages {
			if err := stage(st); err != nil {
				return err
			}
		}
		return nil
	}
}

// NewPipelineState creates a new PipelineState with the provided context and configuration
func NewPipelineState(ctx context.Context, conf *config.AppConfig, client *telegram.Client) *PipelineState {
	return &PipelineState{
		Ctx:    ctx,
		Conf:   conf,
		Client: client,
	}
}

// SetMessages sets the messages in the PipelineState and updates the offset ID
func (ps *PipelineState) SetMessages(msgs []*tg.Message) {
	ps.Messages = msgs
	if len(msgs) > 0 {
		ps.messagesOffsetID = msgs[len(msgs)-1].ID
	} else {
		ps.messagesOffsetID = 0
	}
}
