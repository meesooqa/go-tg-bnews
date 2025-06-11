package proc

import "github.com/gotd/td/tg"

// MessageFilter defines a function type that filters messages based on certain criteria
type MessageFilter func(msg tg.MessageClass) bool

// FilterMessages filters a slice of messages based on the provided filters
func FilterMessages(msgs []tg.MessageClass, filters ...MessageFilter) []tg.MessageClass {
	var out []tg.MessageClass
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
