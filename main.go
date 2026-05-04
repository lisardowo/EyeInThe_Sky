package main

import (
	"EyeInThe_Sky/createConnection"
	"fmt"
	"time"
)

type ServerConnection struct {
	USBAddr string
	IPAddr  string
}

var analysisMode createConnection.TrustLevel = createConnection.Unsecure // TODO HARCORDED VALUE

func main() {


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
	prodServerUSB := ServerConnection{
		USBAddr: "/dev/ttyUSB0",
		IPAddr:  "",
	}
	fmt.Println("--- INITIATING USB PRODUCTION NODE ---")
	createConnection.BootNode(prodServerUSB) // Usar función exportada

	// Test 2: SSH connection of trusted sv
	prodServerSSH := ServerConnection{
		USBAddr: "", // Cambiado a 99 para forzar el fallback a SSH
		IPAddr:  "10.0.0.5",
	}
	fmt.Println("\n--- INITIATING PRODUCTION NODE (SSH FALLBACK) ---")
	createConnection.BootNode(prodServerSSH)

	// Test 3: Unsecure(pentest node) via usb -> worked but do not detect is a untrusted source 
	fmt.Println("\n--- INITIATING UNTRUSTED USB NODE ---")
	pentestNodeUSB := ServerConnection{
		USBAddr: "/dev/ttyUSB99",
		IPAddr:  "",
	}
	createConnection.BootNode(pentestNodeUSB)

	// Test 4: Unsecure(pentest node) via ssh ->not working ip detection stuff || BOOT error works nice tho
	fmt.Println("\n--- INITIATING UNTRUSTED SSH NODE ---")
	pentestNodeSSH := ServerConnection{
		USBAddr: "",
		IPAddr:  "192.168.1.100",
	}
	createConnection.BootNode(pentestNodeSSH)
}

func (server ServerConnection) FetchLatest() (createConnection.TelemetryFrame, error) {
	return createConnection.TelemetryFrame{
		SourceID: "M4C-GENERIC-NODE",
		Level:    createConnection.Secure,
		Data:     "HEARTBEAT_SIGNAL_STABLE",
		Latency:  time.Millisecond * 20,
	}, nil
}

func(server ServerConnection) DetectConnection() (string, bool, error) {
	fmt.Printf("[EYE IN THE SKY] Looking for USB at %s...\n", server.USBAddr)

	usbAvailable, err := createConnection.CheckUSB(server.USBAddr)

	if usbAvailable == int(createConnection.Secure) {
		return "PHYSICAL_USB", true, nil
	} else if usbAvailable == int(createConnection.Unsecure) {
		return "PHYSICAL_USB", false, nil
	}

	fmt.Println(err)
	fmt.Printf("[EYE IN THE SKY] looking for ssh connection at %s..\n", server.IPAddr)

	netStatus := createConnection.CheckNetwork(server.IPAddr)

	if netStatus == int(createConnection.Secure) {
		return "NETWORK_SSH", true, nil
	} else if netStatus == int(createConnection.Unsecure) {
		return "NETWORK_SSH", false, nil
	}

	return "", false, fmt.Errorf("[EYE IN THE SKY] node unreachable: all connections failed")

}

func filterUnsecure(frames []createConnection.TelemetryFrame) []createConnection.TelemetryFrame{
	
	var unsecureFrame []createConnection.TelemetryFrame

	for _,f := range frames{
		if f.Level == createConnection.Unsecure {
			unsecureFrame = append(unsecureFrame,f)
		}
	}
	return unsecureFrame
}

