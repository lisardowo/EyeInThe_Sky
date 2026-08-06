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

//var TrustLevel connection.TrustLevel = connection.Secure // TODO HARCORDED VALUE

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

	bootTime := time.Now()
		initialModel := tui.Model{
 				WhichScreen: tui.HomeScreen, // Default boot time screen, Maybe add memory to keep track of that.. ?
    			Home:        tui.HomeState{
				TrustLevel: connection.Unsecure,
    			Operator:     "debug",
    			VLAN:         0,
    			BootAt:       bootTime,
    			Uptime:       time.Since(bootTime),
   
				},
    			Dash:        tui.DashState{
				TrustLevel:     connection.Unsecure,
    			FocusedPanel: 0,
    			CPUUsage:     67.9,
    			RAMUsage:    69.8,
    			LogsBuffer:   nil,
    			Width:        120 },
    			Width:       120,
	    		Height:      240,
    			LastKey: "NAN",
    			LastAction:  "System Booted Up",

	}
	
	p := tea.NewProgram(initialModel, tea.WithAltScreen(), 
)

    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
	}

}
