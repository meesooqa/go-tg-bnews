package proc

import (
	"github.com/gotd/td/tg"
)

// SkipNoText filters out messages that do not contain any text
func SkipNoText(m *tg.Message) bool {
	return m.String() == ""
}

// KeepLightning filters out messages that do not start with the "⚡️" prefix
func KeepLightning(m *tg.Message) bool {
	if len(m.Entities) == 0 || m.Entities[0].GetOffset() != 0 {
		//return true
		return false
	}
	ent := m.Entities[0]

	s := m.String()
	prefix := s[:ent.GetLength()]
	return prefix != "⚡" && prefix != "⚡️"
}
