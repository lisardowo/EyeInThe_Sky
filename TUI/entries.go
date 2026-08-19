package TUI

import (
	"fmt"
	"strings"
)

type LogCategory int

const (
	CatSystem LogCategory = iota
	CatNetwork
	CatApp
	CatFirewall
)

type LogEntry struct {
	Timestamp	string
	Level 		string
	Message 	string
	Category	LogCategory
}

func (l LogEntry) FormatEntry() string{
	return fmt.Sprintf("[%s] %-5s | %s", l.Timestamp, l.Level, l.Message)
}

func (c LogCategory) IntToStringCategory() string{
	switch c{
		case CatSystem:
			return "SYS"
		case CatNetwork:
			return "NET"
		case CatApp:
			return "APP"
		case CatFirewall:
			return "FW"
		default:
			return "???"
	}
}

func SelectCategory(cmd string) (LogCategory, bool){
	switch strings.ToLower(strings.TrimPrefix(cmd, ":")){ //Lowers, trims the ":" and compares, vim does not do this, it takes the raw command you input, should I do so?
		case "n":
			return CatNetwork, true
		case "s":
			return CatSystem, true
		case "f":
			return CatFirewall, true
		default:
			return 0, false
	}
}

func FilterEntriesCategory(entries []LogEntry, cat LogCategory) []LogEntry{
	out := make([]LogEntry, 0, len(entries))
	for _, e := range entries{
		if e.Category == cat{
			out = append(out,e)
		}
	}
	return out
}