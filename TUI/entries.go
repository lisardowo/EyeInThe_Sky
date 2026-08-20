package TUI

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogCategory int

const (
	CatSystem LogCategory = iota
	CatNetwork
	CatApp
	CatFirewall
	CatProcess
)

const notFound int = -1

type LogEntry struct {
	Timestamp	string
	Level 		string
	Message 	string
	Category	LogCategory
}

func now() string {
	return time.Now().Format("15:04:05")
}

func (l LogEntry) FormatEntry() string{
	return fmt.Sprintf("[%s] %-5s | %s ||||||||||\n", l.Timestamp, l.Level, l.Message)
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

func FormatFiltered(entries []LogEntry, activeFilter *LogCategory) []string{
	if activeFilter != nil{
		entries = FilterEntriesCategory(entries, *activeFilter)
	}
	out := make([]string, 0 , len(entries))
	for _, entry := range entries{
		out = append(out, entry.FormatEntry())
	}

	return out
}

func ReadProcessLogs() ([]LogEntry, error){
	
	procDir, err := os.Open("/proc")
	if err != nil{
		return nil, err
	}
	names, err := procDir.Readdirnames(-1) // -1 forces the function to list every dir it can find
	if err != nil{
		return nil, err
	}

	var entries []LogEntry
	const maxEntries = 20
	for _, name := range names{
		if len(entries) >= maxEntries {
            break
        }
		pid, err := strconv.Atoi(name)
			if err != nil{
				continue // Usually the dir is not a PID but a string name, P.e /proc/meminfo
			}
			entry, ok := readSingleProcess(pid)
			
			if ok {
				entry = FormatFiltered(entry, CatProcess)
				entries = append(entries,entry) //returns a collection/slice of the entries 
			}
	}
	return entries, nil
}

func readSingleProcess(pid int) (LogEntry, bool){
	
	statPath := fmt.Sprintf("/proc/%d/stat", pid) // is this vulnerable to create a pseudo PID or sum??
	data, err := os.ReadFile(statPath)
	if err != nil{
		return LogEntry{}, false // Process can die during the listing and reading
	}
	content := string(data)
	openParen := strings.Index(content, "(") // returns -1 if the substring is not present in the string
	closeParen := strings.LastIndex(content, ")")
	if openParen == notFound || closeParen == notFound || closeParen < openParen{
		return LogEntry{}, false
	}

	command := content[openParen + 1 : closeParen]
	remainingFields := strings.Fields(content[closeParen+1:])
	if len(remainingFields) < 1 {
		return LogEntry{}, false
	}

	state := remainingFields [0]

	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil{
		exePath = "?" // Kernel threads/unauthorized access wont resolve
	}

	level := "SAFE"

	if strings.HasPrefix(exePath, "/tmp") || strings.HasPrefix(exePath, "/dev/shm"){
		level = "WARN" // If is being ran from a not std location, assume is not safe
	}

	return LogEntry{
		Timestamp:  now(),
		Level:		level,
		Category:	CatProcess,
		Message:	fmt.Sprintf("pid:%-6d %-16s state:%-2s exe:%s", pid, command, state, exePath),
	}, true

}