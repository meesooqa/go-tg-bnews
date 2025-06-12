package proc

import (
	"strings"

	"github.com/gotd/td/tg"
)

// SkipNoText filters out messages that do not contain any text
func SkipNoText(m *tg.Message) bool {
	return m.Message == ""
}

// KeepLightning filters out messages that do not start with the "⚡️" prefix
func KeepLightning(m *tg.Message) bool {
	if m == nil {
		return true
	}
	return !strings.HasPrefix(m.Message, "⚡️")
}
