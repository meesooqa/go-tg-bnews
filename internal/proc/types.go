package proc

import (
	"context"
	"sync"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	mytg "github.com/meesooqa/go-tg-bnews/internal/telegram"
)

// PipelineState holds the state of the processing pipeline
type PipelineState struct {
	Ctx         context.Context
	Client      *telegram.Client
	Service     *mytg.Service
	ChannelFrom *tg.Channel
	ChannelTo   *tg.Channel
	Messages    []*tg.Message

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
