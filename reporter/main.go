package reporter

import (

	"bufio"
	"encoding/json"
	"fmt"
	"go/scanner"
	"os"
	"strings"
	"time"
	"EyeInTheSky/TUI"
	"EyeInTheSky/helpers"
	"EyeInTheSky/sysinfo"

	"go.bug.st/serial"
)

const (
	BaudRate  = 115200//TODO harcoded to std, use unix library to get dinamically
	MAGICSYN  = "SYN:EYE_IN_THE_SKY" 
	MAGICACK  = "ACK:EYE_REPORTER_ONLINE"
	MAGICHALT = "HALT"

)

type report struct {
	allEntries 	[]helpers.LogEntry
	cpu			sysinfo.CPUSample
	ram			float64
}

func resolvePort() string{
	if len(os.Args) > 1{
		return os.Args[1]
	}

	if envPort := os.Getenv("AGENT_PORT"); envPort != "" {
		return envPort
	}

	ports, err := serial.GetPortsList()
	if err == nil {
		for _, port := range ports {
			if strings.Contains(port, "ttyUSB") || strings.Contains(port, "ttyACM"){
				return port
			}
		}
		if len(ports) > 0 {
			return ports[0]
		}
	}

	return ""

}

func main() { 

	mode := &serial.Mode{
		BaudRate: BaudRate,
		DataBits: 8,
		Parity: serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	for {
		portPath := resolvePort()
		if portPath == ""{
			time.Sleep(2 * time.Second)
			continue
		}

		port, err := serial.Open(portPath, mode)

		if err != nil {
			time.Sleep(2 * time.Second)
		}

		fmt.Printf("[DAEMON] Listening at %s. Waiting for handshake..\n", portPath)
		scanner := bufio.NewScanner(port)

		activeSession := false
		for scanner.Scan(){
			if strings.TrimSpace(scanner.Text()) == MAGICSYN{
				_, _ = port.Write([]byte(MAGICACK + "\n"))
				activeSession = true
				break
			}
		}
		if activeSession {
			fmt.Println("[DAEMON] Handshake successful ")
			streamMetrics(port)
		}

		port.Close()
		time.Sleep(1 * time.Second)

	}

}

func streamMetrics(port serial.Port){
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	sampleNow, err := sysinfo.GetCPUSample()
	if err != nil {
		prevSample = sysinfo.CPUSample{}
	}


	for range ticker.C {

		var all []helpers.LogEntry

		if e, err := helpers.ReadProcessLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadTCPLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadUDPLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadModulesLogs(); err == nil {
			all = append(all, e...)
		}
		if e, err := helpers.ReadDiskstatsLogs(); err == nil {
			all = append(all, e...)
		}
		ramU, _ :=  sysinfo.GetRamUsage()
		payload := report{
			
			allEntries: all,
			cpu:  	sysinfo.CPUSample{
						Idle:  sampleNow.Idle,
						Total: sampleNow.Total
					},
			ram:  		ramU,
		}
	}
}