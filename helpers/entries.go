package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//TODO move helper function to helpers module and leave in here just the logic to read and process final entries
type LogCategory int
const diskStateFields int = 14

const (
	CatSystem LogCategory = iota
	CatNetwork
	CatApp
	CatFirewall
	CatProcess
	CatKernel
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
	return fmt.Sprintf("[%s] -- %s -- : %s", l.Timestamp, l.Level, l.Message)
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
		case CatProcess:
			return "PROC"
		case CatKernel:
			return "KRNL"
		
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
		case "a":
			return CatApp, true
		case "p":
			return CatProcess, true
		case "k":
			return CatKernel, true
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
	
	for _, name := range names{
		
		pid, err := strconv.Atoi(name)
			if err != nil{
				continue // Usually the dir is not a PID but a string name, P.e /proc/meminfo
			}
			entry, ok := readSingleProcess(pid)
			
			if ok {
				
				entries = append(entries,entry )
				
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
		Message:	fmt.Sprintf("pid:%-6d %-16s state:%-2s exe:%s\n", pid, command, state, exePath), // message is constructed in two moments
	}, true

}

func readNetworkTable(path string, protocol string)([]LogEntry, error){
	fd, err := os.Open(path) // function meant for both TCP and udp
	if err != nil{
		return nil, err
	}
	defer fd.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(fd)
	line := 0
	for scanner.Scan(){
		line ++
		if line == 1 {
			continue // header, ignore 
		}
		//2664A8C0:D022 IP:PORT (?
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue //TODO if not Valid hex len ignore, should throw err? 
		}

		localAddress := HexToIPport(fields[1])
		remoteAddress := HexToIPport(fields[2])
		state := GetState(fields[3])
		if state == "NOT VALID"{
			return entries, scanner.Err()
		}
		entries = append(entries, LogEntry{
			Timestamp: now(),
			Level: "INFO",
			Category: CatNetwork,
			Message: fmt.Sprintf("[%s] %s -> %s : State{%s}\n", protocol, localAddress,remoteAddress, state),
		})

	}

	return entries, nil 

}

func ReadTCPLogs() ([]LogEntry, error){
	return readNetworkTable("/proc/net/tcp", "TCP")
}


func ReadUDPLogs() ([]LogEntry, error){
	return readNetworkTable("/proc/net/udp", "UDP")
}

func readModulesLogs()([]LogEntry, error){
	fd, err := os.Open("/proc/modules")
	if err != nil{
		return nil, err
	}

	defer fd.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(fd)
	for scanner.Scan(){
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3{
			continue
		}
		name := fields[0]
		size := fields[1]
		useCount := fields[2]

		entries = append(entries, LogEntry{
			Timestamp: now(),
			Level: "INFO",
			Category: CatKernel,
			Message: fmt.Sprintf("%-20s size:%-8s used:%s", name, size, useCount)})
	}
	return entries, scanner.Err()
}

func ReadDiskstatsLogs()([]LogEntry, error){
	fd, err := os.Open("/proc/diskstats")
	if err != nil{
		return nil, err
	}

	defer fd.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < diskStateFields{
			continue
		}
		device := fields[2]
		sectorsRead := fields[5]
		sectorsWritten := fields[9]

		if strings.HasPrefix(device, "loop") || strings.HasPrefix(device, "ram"){
			continue //ignore loop/ram devices since dont have valuable information
		}
		entries = append(entries, LogEntry{
			Timestamp: now(),
			Level: "INFO",
			Category: CatSystem,
			Message: fmt.Sprintf("%-8s sectors_read:%-10s sectors_written:%s", device, sectorsRead, sectorsWritten),
		})
	}
	return entries, scanner.Err()
}