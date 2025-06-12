package proc

import (
	"testing"

	"github.com/gotd/td/tg"
	"github.com/stretchr/testify/assert"
)

func TestSkipNoText(t *testing.T) {
	tests := []struct {
		name     string
		message  *tg.Message
		expected bool
	}{
		{
			name: "Empty message",
			message: &tg.Message{
				Message: "",
			},
			expected: true,
		},
		{
			name: "Non-empty message",
			message: &tg.Message{
				Message: "Hello",
			},
			expected: false,
		},
		{
			name: "Whitespace-only message",
			message: &tg.Message{
				Message: "   ",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SkipNoText(tt.message))
		})
	}
}

func TestKeepLightning(t *testing.T) {
	tests := []struct {
		name string
		msg  *tg.Message
		want bool
	}{
		{
			name: "nil message",
			msg:  nil,
			want: true,
		},
		{
			name: "empty text",
			msg:  &tg.Message{Message: ""},
			want: true,
		},
		{
			name: "with prefix",
			msg:  &tg.Message{Message: "⚡️Hello, world!"},
			want: false,
		},
		{
			name: "without prefix",
			msg:  &tg.Message{Message: "Hello ⚡️ world!"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeepLightning(tt.msg)
			if got != tt.want {
				t.Errorf("HasLightningPrefix(%v) = %v; want %v", tt.msg, got, tt.want)
			}
		})
	}
}
