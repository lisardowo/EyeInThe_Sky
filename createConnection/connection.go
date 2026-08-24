package createconnection

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.bug.st/serial"
)

type TrustLevel int

const (
	Secure   TrustLevel = iota //0s
	Unsecure
)

func (t TrustLevel) TrustToString() string{
    switch t {
        
        case Secure:
            return "Secure"
        case Unsecure:
            return "Unsecure"
        default: 
            return "Unknown"

    }

}


type RemoteFrame struct {
    CPU         float64 `json:"cpu"`
    RAM         float64 `json:"ram"`
    Log         string  `json:"log"`
    Category    int      `json:"category"`
}

type USBConnection struct {
    PortPath string
    BaudRate int
    Port     serial.Port
    Scanner  *bufio.Scanner
}

func NewUSBConnection(portPath string, baudRate int) *USBConnection {
    if portPath == "" {
        portPath = DetectUSB()
    }
    if baudRate == 0 {
        baudRate = 115200// returns std baudRate 4 lunux
    }

    return &USBConnection{
        PortPath: portPath,
        BaudRate: baudRate,
    }
}

func DetectUSB() string {
    ports, err := serial.GetPortsList()
    if err != nil || len(ports) == 0 {
        return ""
    }

    for _, port := range ports {
        if strings.Contains(port, "ttyUSB") || strings.Contains(port, "ttyACM") {
            return port
        }
    }
    return ports[0]
}

func (connection *USBConnection) Connect() error{

    if connection.PortPath == ""{
        return errors.New("Unable to detect an USB/serial device")
    }

    if _, err := os.Stat(connection.PortPath); os.IsNotExist(err) {
        return fmt.Errorf("%s port does not exist", connection.PortPath)
    }

    mode := &serial.Mode {

        BaudRate: connection.BaudRate,
        DataBits: 8,
        Parity: serial.NoParity,
        StopBits: serial.OneStopBit,
    
    }

    port, err := serial.Open(connection.PortPath, mode)

    if err != nil {
        return fmt.Errorf("Unable to open: %s", err)
    }

    if err := port.SetReadTimeout(2 * time.Second); err != nil {
        port.Close()
        return err
    }

    connection.Port = port
    connection.Scanner = bufio.NewScanner(port)
    return nil

}