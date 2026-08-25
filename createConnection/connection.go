package createconnection

import (
	"bufio"
	"encoding/json"
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

func NewUSBConnection() *USBConnection {
    
    portPath := DetectUSB()
    baudRate := getBaud()

    if baudRate == 0 {
        baudRate = 115200// returns std baudRate 4 lunux
    }

    return &USBConnection{
        PortPath: portPath,
        BaudRate: baudRate,
    }
}

func getBaud() int{

    return  115200 //TODO std linux baud rate, not actual tho.. Need to implement a func to dinamic get that 

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

func (connection *USBConnection) Close() error {
   if connection.Port != nil{
        return connection.Port.Close()
   }
   return nil
}

func (connection *USBConnection) ReadNextFrame() (RemoteFrame,error){

    if connection.Port != nil{
        return RemoteFrame{}, fmt.Errorf("Port was not initialized")
    }

    if connection.Scanner.Scan(){
        line := strings.TrimSpace(connection.Scanner.Text())
        if len(line) == 0 {
            return RemoteFrame{}, fmt.Errorf("Empty frame")
        }

        var frame RemoteFrame

        if err := json.Unmarshal([]byte(line), &frame); err != nil{
            return RemoteFrame{}, fmt.Errorf("Error while parsing frame: %w", err)
        }
        return frame,nil
    }

    if err := connection.Scanner.Err(); err != nil{
        return RemoteFrame{}, err
    }


    return RemoteFrame{}, fmt.Errorf("USB timeout")
}

func (connection *USBConnection) PerformHandshake() error {
    if connection.Port == nil {
        return errors.New("Port not initialized")
    }

    if _, err := connection.Port.Write([]byte("SYN:EYE_IN_THE_SKY\n")); err != nil{
        return fmt.Errorf("[SYN] Error performing handshake in %v", err)
    }

    if connection.Scanner.Scan(){
        resp := strings.TrimSpace(connection.Scanner.Text())

        if resp == "ACK:EYE_REPORTER_ONLINE"{
            return nil
        }
        return fmt.Errorf("Invalid answer from %s", resp)
    }

    return fmt.Errorf("timeout reached waiting for the demon")

}