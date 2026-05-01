package main

import (
	"fmt"
	"time"
)

type TrustLevel int

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

func main(){
	labFrame := TelemetryFrame{
		SourceID: "CompletelyRealTestFrame",
		Level: Unsecure,
		Data: "inbound connection detected on sum port ",
		Latency: time.Millisecond * 15,
	}

	fmt.Printf("---EYE IN THE SKY --- \n")
	fmt.Printf("Source: %s\n DATA: %s\n", labFrame.SourceID, labFrame.Data)
}