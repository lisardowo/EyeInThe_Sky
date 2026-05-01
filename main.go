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

var analysisMode TrustLevel = Unsecure //TODO harcorded value for testing only, the analysisMode restrains the capabilities
//of communication, in a secure source, communication can be lighter, in a unsecure communication is way more restricted, besides other UI/UX changes
//Analysis Mode is supposed to be obtained via a "handshake" when both devices are connected, before anything starts

func main(){


	currentLabFrame := TelemetryFrame{
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
	}
	
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