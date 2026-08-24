package createconnection

import (
	"bufio"

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
        //implement a function to check the usb connection
    }
    if baudRate == 0 {
        baudRate = 115200// returns std baudRate 4 lunux
    }

    return &USBConnection{
        PortPath: portPath,
        BaudRate: baudRate,
    }
}
