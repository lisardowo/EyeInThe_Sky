package TEST

import (
	connection "EyeInThe_Sky/createConnection"
	"fmt"
	"strings"

	//"strings"
	"time"
)

type ServerConnection struct {
	USBAddr string
	IPAddr  string
}

func (server ServerConnection) FetchLatest() (connection.TelemetryFrame, error) {
	return connection.TelemetryFrame{
		SourceID: "M4C-GENERIC-NODE",
		Level:    connection.Secure,
		Data:     "HEARTBEAT_SIGNAL_STABLE",
		Latency:  time.Millisecond * 20,
	}, nil
}

func(server ServerConnection) DetectConnection() (string, bool, error) {
	fmt.Printf("[EYE IN THE SKY] Looking for USB at %s...\n", server.USBAddr)

	usbAvailable, err := connection.CheckUSB(server.USBAddr)

	if usbAvailable == int(connection.Secure) {
		return "PHYSICAL_USB", true, nil
	} else if usbAvailable == int(connection.Unsecure) {
		return "PHYSICAL_USB", false, nil
	}

	fmt.Println(err)
	fmt.Printf("[EYE IN THE SKY] looking for ssh connection at %s..\n", server.IPAddr)

	netStatus := connection.CheckNetwork(server.IPAddr)

	if netStatus == int(connection.Secure) {
		return "NETWORK_SSH", true, nil
	} else if netStatus == int(connection.Unsecure) {
		return "NETWORK_SSH", false, nil
	}

	return "", false, fmt.Errorf("[EYE IN THE SKY] node unreachable: all connections failed")

}

func filterUnsecure(frames []connection.TelemetryFrame) []connection.TelemetryFrame{
	
	var unsecureFrame []connection.TelemetryFrame

	for _,f := range frames{
		if f.Level == connection.Unsecure {
			unsecureFrame = append(unsecureFrame,f)
		}
	}
	return unsecureFrame
}

func parseTrustLevel(value string) connection.TrustLevel {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "secure", "s", "trusted":
		return connection.Secure
	default:
		return connection.Unsecure
	}
}

