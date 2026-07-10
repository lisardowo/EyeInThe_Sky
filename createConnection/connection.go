package connection

import (
	"fmt"
	"strings"
	"time"
)

type TrustLevel int

const (
	Secure   TrustLevel = iota
	Unsecure
)

func (t TrustLevel) String() string {
    switch t {
    case Secure:
        return "Secure"
    case Unsecure:
        return "Unsecure"
    default:
        return "Unknown"
    }
}

type TelemetryFrame struct {
	SourceID string
	Level    TrustLevel
	Data     string
	Latency  time.Duration
}

type SourceNode interface {
	DetectConnection() (string, bool, error)
	FetchLatest() (TelemetryFrame, error)
}

func CheckUSB(addr string) (int, error) { 
    if addr == "/dev/ttyUSB0" {
        return int(Secure), nil
    } else if strings.HasPrefix(addr, "/dev/ttyUSB") {
        return int(Unsecure), nil
    }
    return -1, fmt.Errorf("[EYE IN THE SKY] USB unavailable. Attempting SSH fallback")
}

func CheckNetwork(ip string) int { 
    if strings.HasPrefix(ip, "10.0.0") {
        return int(Secure)
    }
    return int(Unsecure)
}

func BootNode(node SourceNode) { 
	transport, isSecure, err := node.DetectConnection()

    if err != nil {
        fmt.Printf("[EYE IN THE SKY] BOOT ERROR: %v\n", err)
        return
    }

    fmt.Printf("--------------------------------------\n")
    fmt.Printf("[EYE IN THE SKY] NODE ONLINE | Transport: %s\n", transport)

    if !isSecure {
        fmt.Println("[EYE IN THE SKY] <WARNING> Running in Untrusted Source")
    } else {
        fmt.Println("[EYE IN THE SKY] SECURE CONNECTION ESTABLISHED")
    }

    frame, _ := node.FetchLatest()
    fmt.Printf("[EYE IN THE SKY] DATA: [%s] %s\n", frame.SourceID, frame.Data)
    fmt.Printf("--------------------------------------\n")
}