package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

const ipv4Size int = 4

func HexToIPport(field string) string{
	parts := strings.Split(field, ":")
	if len(parts) != 2{
		return field
	}
	ipHex, portHex := parts[0], parts[1]
	if len(ipHex) != 8{
		return field // Ignoring IPV6 rn
	}
	var ipFragments [ipv4Size]byte
	for i  := 0; i < ipv4Size ; i++{
		decimal, err := strconv.ParseUint(ipHex[i*2:i*2+2], 16, 8)
		if err != nil{
			return field
		}
		
		ipFragments[i] = byte(decimal)
	}
	port, err := strconv.ParseUint(portHex, 16, 16)
	
	if err != nil{
		return field
	}

	return fmt.Sprintf("%d.%d.%d.%d:%d", ipFragments[0], ipFragments[1],ipFragments[2],ipFragments[3], port)

}

func GetState(statehex string) string{
	switch statehex{
		case "01":
			return "ESTABLISHED"
		case "02":
			return "SYN_SENT"
		case "03":
			return "SYN_RECV"
		case "04":
			return  "FIN_WAIT1"
		case "05":
			return	"FIN_WAIT2"
		case "06":
			return "TIME_WAIT"
		case "07":
			return "CLOSE"
		case "08":
			return "CLOSE_WAIT"
		case "09":
			return	"LAST_ACK"
		case "0A":
			return	"LISTEN"
		case "0B":
			return	"CLOSING"
		case "0C":
			return  "NEW_SYN_RECV"
		default:
			return "NOT VALID" //TODO Should this return a warning alrt?
	}
}