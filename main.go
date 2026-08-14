package main

import (
	tui "EyeInThe_Sky/TUI"
	connection "EyeInThe_Sky/createConnection"
	"flag"
	"fmt"
	"os"
	"os/user"
	"time"

	//"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerConnection struct {
	USBAddr string
	IPAddr  string
}


func main() {


	//modeFlag := flag.String("mode", TrustLevel.String(), "analysis mode: secure or unsecure")
	operatorFlag := flag.String("operator", "", "operator name shown in the TUI")
	//vlanFlag := flag.Int("vlan", 10, "VLAN identifier shown in the TUI")
	flag.Parse()

	operator := *operatorFlag
	if operator == "" {
		if currentUser, err := user.Current(); err == nil && currentUser.Username != "" {
			operator = currentUser.Username
		} else {
			operator = "operator"
		} 	
	}

	// Creates a channel and sends it as struct to the home render 

	BootUp := time.Now()
	var UptimeTimer chan time.Duration
	UptimeTimer = make(chan time.Duration)

	go tui.UptimeHandle(BootUp, UptimeTimer)
	ringBuffer := tui.NewBuffer[any](50)
	initialModel := tui.Model{
		
 				WhichScreen: tui.HomeScreen, // Default boot time screen, Maybe add memory to keep track of that.. ?
    			Width:       120,
	    		Height:      240,
    			LastKey: "NAN",
    			LastAction:  "System Booted Up",
				TrustLevel:     connection.Secure,
				
				Home:        tui.HomeState{
    			Operator:     "debug",
    			VLAN:         0,
    			Uptime:      UptimeTimer,

				},

    			Dash:        tui.DashState{
    			FocusedPanel: 0,
    			CPUUsage:     67.9,
    			RAMUsage:    69.8,
    			LogsBuffer:   *ringBuffer,/* tui.RingBuffer[any]{
					Buffer: nil,
					Size: 50,
					Head: 0,
					Count: 0,
				} */ 
    			},

	}
	
	p := tea.NewProgram(initialModel, tea.WithAltScreen(), 
)

    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
	}

}
