package sysinfo

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)


const(

	metricIdentifier = 0
	userTime = 1
	idleTime = 4	
)

type CPUSample struct {
	Idle uint64
	Total uint64
}

func GetCPUSample()(CPUSample, error){
	file, err := os.Open("/proc/stat")
	 
	if err != nil {
		return CPUSample{}, err //TODO add this logs of the tool 
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	if scanner.Scan(){
		fields := strings.Fields(scanner.Text())
		
		if len(fields) >= 5 && fields[metricIdentifier] == "cpu" {
			var total uint64
			for _, val := range fields[userTime:]{
				n, _ := strconv.ParseUint(val, 10, 64)
				total += n
			}
			idle, _ := strconv.ParseUint(fields[idleTime], 10 , 64)
			return CPUSample{Idle: idle, Total: total}, nil
		}
	}
	return CPUSample{}, nil
}