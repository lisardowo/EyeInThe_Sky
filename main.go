package main

import (
	tui "EyeInThe_Sky/TUI"
	connection "EyeInThe_Sky/createConnection"
	"flag"
	"fmt"
	"os"
	"os/user"

	//"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ServerConnection struct {
	USBAddr string
	IPAddr  string
}

var analysisMode connection.TrustLevel = connection.Unsecure // TODO HARCORDED VALUE

func main() {

	fmt.Printf("sexo")

	//modeFlag := flag.String("mode", analysisMode.String(), "analysis mode: secure or unsecure")
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

	//parsedMode := parseTrustLevel(*modeFlag)


	/*currentLabFrame := TelemetryFrame{
		SourceID: "CompletelyRealTestFrame",
		Level: Unsecure,
		Data: "inbound connection detected on sum port ",
		Latency: time.Millisecond * 15,
	}

	fmt.Printf("---EYE IN THE SKY --- \n")
	fmt.Printf("Source: %s\n DATA: %s\n", currentLabFrame.SourceID, currentLabFrame.Data)
	if analysisMode == Unsecure{
		fmt.Println("Analyzing payload from an untrusted source")
		
		input := []TelemetryFrame{currentLabFrame}
		unsecurePayloads := filterUnsecure(input)
		fmt.Printf("Filtered Frames Count: %d\n", len(unsecurePayloads))

	}else if analysisMode == Secure{
		fmt.Println("Analyzing payload from a trusted source")
	}*/
	//TODO testing connection mock functions

	// Test 1: USB connection of trusted sv
	/* TODO testing OPEN a tui
	prodServerUSB := ServerConnection{
		USBAddr: "/dev/ttyUSB0",
		IPAddr:  "",
	}
	fmt.Println("--- INITIATING USB PRODUCTION NODE ---")
	connection.BootNode(prodServerUSB) // Usar función exportada

	// Test 2: SSH connection of trusted sv
	prodServerSSH := ServerConnection{
		USBAddr: "", // Cambiado a 99 para forzar el fallback a SSH
		IPAddr:  "10.0.0.5",
	}
	fmt.Println("\n--- INITIATING PRODUCTION NODE (SSH FALLBACK) ---")
	connection.BootNode(prodServerSSH)

	// Test 3: Unsecure(pentest node) via usb -> worked but do not detect is a untrusted source 
	fmt.Println("\n--- INITIATING UNTRUSTED USB NODE ---")
	pentestNodeUSB := ServerConnection{
		USBAddr: "/dev/ttyUSB99",
		IPAddr:  "",
	}
	connection.BootNode(pentestNodeUSB)

	// Test 4: Unsecure(pentest node) via ssh ->not working ip detection stuff || BOOT error works nice tho
	fmt.Println("\n--- INITIATING UNTRUSTED SSH NODE ---")
	pentestNodeSSH := ServerConnection{
		USBAddr: "",
		IPAddr:  "192.168.1.100",
	}
	connection.BootNode(pentestNodeSSH) */
		
		initialModel := tui.Model{
 			CurrentMode: 0,
    		Home:        tui.HomeState{},
    		Dash:        tui.DashState{},
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
