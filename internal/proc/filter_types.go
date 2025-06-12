package proc

import "github.com/gotd/td/tg"

// MessageFilter defines a function type that filters messages based on certain criteria
type MessageFilter func(msg *tg.Message) bool

// FilterMessages filters a slice of messages based on the provided filters
func FilterMessages(msgs []*tg.Message, filters ...MessageFilter) []*tg.Message {
	var out []*tg.Message
	for _, m := range msgs {
		skip := false
		for _, f := range filters {
			if f(m) {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, m)
		}
	}
	return out
}
