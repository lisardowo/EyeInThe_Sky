package main

import (
	"fmt"
	"strings"
	"time"
)

type TrustLevel int

type SourceNode interface{

	DetectConnection() (string, error)
	//HandkShake() (TrustLevel, error)
	Authenticate() (TrustLevel, error)	
	FetchLatest() (TelemetryFrame, error)

}

type ServerConnection struct{
	USBAddr string
	IPAddr	string
}

type ProxmoxNode struct{
	IP string
}

//func (proxNode ProxmoxNode) Handshake() (TrustLevel, error){}

const (
	
	Secure TrustLevel = iota //iota is a number generator -> this will be 0
	Unsecure // and this 1 (kind of stupid tho)

)

type TelemetryFrame struct {
	SourceID string
	Level	TrustLevel
	Data	string
	Latency	time.Duration
}

var analysisMode TrustLevel = Unsecure //TODO harcorded value for testing only, the analysisMode restrains the capabilities
//of communication, in a secure source, communication can be lighter, in a unsecure communication is way more restricted, besides other UI/UX changes
//Analysis Mode is supposed to be obtained via a "handshake" when both devices are connected, before anything starts



func main(){


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
	BootNode(prodServerUSB)

	// Test 2: SSH connection of trusted sv
	prodServerSSH := ServerConnection{
		USBAddr: "/dev/ttyUSB99", // Cambiado a 99 para forzar el fallback a SSH
		IPAddr:  "10.0.0.5",
	}
	fmt.Println("\n--- INITIATING PRODUCTION NODE (SSH FALLBACK) ---")
	BootNode(prodServerSSH)

	// Test 3: Unsecure(pentest node) via usb -> worked but do not detect is a untrusted source 
	fmt.Println("\n--- INITIATING UNTRUSTED NODE ---")
	pentestNodeUSB := ServerConnection{
		USBAddr: "/dev/ttyUSB0",
		IPAddr:  "",
	}
	BootNode(pentestNodeUSB)

	// Test 4: Unsecure(pentest node) via ssh ->not working ip detection stuff || BOOT error works nice tho
	fmt.Println("\n--- INITIATING UNTRUSTED NODE ---")
	pentestNodeSSH := ServerConnection{
		USBAddr: "/dev/ttyUSB1",
		IPAddr:  "192.168.1.100",
	}
	BootNode(pentestNodeSSH)
}

func (server ServerConnection) FetchLatest() (TelemetryFrame, error) {
	return TelemetryFrame{
		SourceID: "M4C-GENERIC-NODE",
		Level:    Secure,
		Data:     "HEARTBEAT_SIGNAL_STABLE",
		Latency:  time.Millisecond * 20,
	}, nil
}

func(server ServerConnection) DetectConnection() (string,error) {
	fmt.Printf("[EYE IN THE SKY] Looking for USB at %s...\n", server.USBAddr)
	usbAvailable := checkUSB(server.USBAddr)

	if usbAvailable{
		return "PHYSICAL_USB", nil
	}

	fmt.Printf("[EYE IN THE SKY] USB unavailable. Attempting SSH fallback to %s...\n", server.IPAddr)
	if checkNetwork(server.IPAddr){
		return "NETWORK_SSH", nil
	}

	return "", fmt.Errorf("node unreachable: all conections failed")

}

func filterUnsecure(frames []TelemetryFrame) []TelemetryFrame{
	
	var unsecureFrame []TelemetryFrame

	for _,f := range frames{
		if f.Level == Unsecure {
			unsecureFrame = append(unsecureFrame,f)
		}
	}
	return unsecureFrame
}

func (server ServerConnection) Authenticate() (TrustLevel, error){
 
    nodeID := server.USBAddr + server.IPAddr

   
    if strings.Contains(nodeID, "192.168") || strings.Contains(nodeID, "VULN") {
        return Unsecure, nil
    }

    return Secure, nil
}


func checkUSB(addr string) bool {
	if addr == "/dev/ttyUSB0"{
		return true
	}
	return false
}

func checkNetwork(ip string) bool {
	if strings.HasPrefix(ip, "10.0.0"){
		return true
	}
	return false
}

func BootNode(node SourceNode) {
    
    transport, err := node.DetectConnection()
    if err != nil {
    
        fmt.Printf("[EYE IN THE SKY] BOOT ERROR: %v\n", err)
        return
    }

    level, err := node.Authenticate()
    if err != nil {
        fmt.Printf("[EYE IN THE SKY]  AUTH ERROR: %v\n", err)
        return
    }

    analysisMode = level

    fmt.Printf("--------------------------------------\n")
    fmt.Printf("[EYE IN THE SKY] NODE ONLINE | Transport: %s\n", transport)
    
    if analysisMode == Unsecure {
        fmt.Println("[EYE IN THE SKY] <WARNING> Running in Untrusted Source")
    } else {
        fmt.Println("[EYE IN THE SKY] SECURE CONNECTION ESTABLISHED")
    }
    
    
    frame, _ := node.FetchLatest()
    fmt.Printf("[EYE IN THE SKY] DATA: [%s] %s\n", frame.SourceID, frame.Data)
    fmt.Printf("--------------------------------------\n")
}