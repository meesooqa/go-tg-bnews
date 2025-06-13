package proc

import (
	"context"
	"errors"
	"testing"

	"github.com/gotd/td/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterMessages(t *testing.T) {
	msg1 := &tg.Message{ID: 1}
	msg2 := &tg.Message{ID: 2}
	msg3 := &tg.Message{ID: 3}
	messages := []*tg.Message{msg1, msg2, msg3}

	// Filter to exclude even IDs
	evenFilter := func(m *tg.Message) bool {
		return m.ID%2 == 0
	}

	// Filter to exclude messages with ID > 2
	idOver2Filter := func(m *tg.Message) bool {
		return m.ID > 2
	}

	t.Run("AllFiltersApplied", func(t *testing.T) {
		result := FilterMessages(messages, evenFilter, idOver2Filter)
		require.Equal(t, []*tg.Message{msg1}, result)
	})

	t.Run("NoFilters", func(t *testing.T) {
		result := FilterMessages(messages)
		require.Equal(t, messages, result)
	})

	t.Run("NoMessages", func(t *testing.T) {
		var empty []*tg.Message
		result := FilterMessages(empty, evenFilter)
		require.Empty(t, result)
	})

	t.Run("SomeMessagesFiltered", func(t *testing.T) {
		singleFilter := func(m *tg.Message) bool { return m.ID == 2 }
		result := FilterMessages(messages, singleFilter)
		require.Equal(t, []*tg.Message{msg1, msg3}, result)
	})
}

func TestChain(t *testing.T) {
	initialMessages := []*tg.Message{{ID: 1}, {ID: 2}, {ID: 3}}
	state := &PipelineState{
		Ctx:      context.Background(),
		Messages: initialMessages,
	}

	t.Run("ExecutionOrder", func(t *testing.T) {
		callOrder := []int{}
		p1 := Processor(func(_ *PipelineState) error {
			callOrder = append(callOrder, 1)
			return nil
		})
		p2 := Processor(func(_ *PipelineState) error {
			callOrder = append(callOrder, 2)
			return nil
		})

		chained := Chain(p1, p2)
		err := chained(state)
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2}, callOrder)
	})

	t.Run("ErrorPropagation", func(t *testing.T) {
		testError := errors.New("test error")
		pErr := Processor(func(_ *PipelineState) error {
			return testError
		})
		pNeverCalled := Processor(func(_ *PipelineState) error {
			t.Fatal("This processor should not be called")
			return nil
		})

		chained := Chain(pErr, pNeverCalled)
		err := chained(state)
		require.Error(t, err)
		assert.Equal(t, testError, err)
	})

	t.Run("StateModification", func(t *testing.T) {
		pAddMessage := Processor(func(s *PipelineState) error {
			s.Messages = append(s.Messages, &tg.Message{ID: 4})
			return nil
		})
		pVerify := Processor(func(s *PipelineState) error {
			assert.Len(t, s.Messages, 4)
			assert.Equal(t, 4, s.Messages[3].ID)
			return nil
		})

		chained := Chain(pAddMessage, pVerify)
		err := chained(state)
		require.NoError(t, err)
	})
}
