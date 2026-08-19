package TUI

import (
	"fmt"
)


type LogEntry struct {
	Timestamp	string
	Level 		string
	Message 	string
}

func (l LogEntry) FormatEntry() string{
	return fmt.Sprintf("[%s] %-5s | %s", l.Timestamp, l.Level, l.Message)
}
