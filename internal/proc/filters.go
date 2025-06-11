package proc

import "github.com/gotd/td/tg"

// SkipNoText filters out messages that do not contain any text
func SkipNoText(m tg.MessageClass) bool {
	return m.String() == ""
}
