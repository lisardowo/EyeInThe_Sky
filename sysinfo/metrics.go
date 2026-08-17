package sysinfo

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

/*

Table of values for os.Open("/proc/stat") (P.e [cpu  237391 183 74700 5025332 22923 14683 9828 0 0 0])
												 0      1    2    3      4    ...
| Index  |Type of time| Description                                                  |
|--------|------------|--------------------------------------------------------------|
| [0]    | cpu        | CPU identifier ("cpu", "cpu0", "cpu1")					     |
| [1]    | User       | Normal programs in user space                   		     |
| [2]    | Nice       | Low-priority user space programs (nice > 0)		             |
| [3]    | System     | Programs executed in kernel space program	                 |
| [4]    | Idle       | CPU idle time								                 |
| [5]    | IOWAIT     | Wating for the disk								           	 |
| [6]    | IRQ        | Hardware interrupts (physical)       	                     |
| [7]    | SOFTIRQ    | Software interrupts (deferred)		                       	 |
| [8]    | Steal      | Stolen time by Hypervisor (important in web configurations(?)|
| [9]    | Guest      | Executing a VM/container					                 |
| [10]   | Guest_Nice | Executing VM with low priority 					             |
*/

const(
	metricIdentifier = 0
	userTime = 1
	idleTime = 4	
	memoryParameter = 1
)

type CPUSample struct {
	Idle uint64
	Total uint64
}
/* TODO corrections for getcpusample
func GetCPUSample() (CPUSample, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return CPUSample{}, fmt.Errorf("failed to open /proc/stat: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		
		// /proc/stat tiene mínimo 5 campos esenciales, pero lo normal son 11 (ID + 10 métricas)
		if len(fields) >= 5 && fields[metricIdentifier] == "cpu" {
			var total uint64
			
			// Para evitar dobles sumas con Guest/GuestNice, sumamos explícitamente los primeros 8 contadores numéricos
			// Desde userTime (1) hasta Steal (8). Ajusta los índices según tus constantes.
			maxFields := len(fields)
			if maxFields > 9 { 
				maxFields = 9 // Limitamos el bucle para no procesar Guest y Guest_Nice si existen
			}

			for i := userTime; i < maxFields; i++ {
				n, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					return CPUSample{}, fmt.Errorf("error parsing metric at index %d: %w", i, err)
				}
				total += n
			}

			// Extraemos el valor Idle de forma segura
			idle, err := strconv.ParseUint(fields[idleTime], 10, 64)
			if err != nil {
				return CPUSample{}, fmt.Errorf("error parsing idle time: %w", err)
			}

			return CPUSample{Idle: idle, Total: total}, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return CPUSample{}, fmt.Errorf("scanner error: %w", err)
	}

	// Si llegamos aquí, la estructura de /proc/stat no era la esperada
	return CPUSample{}, fmt.Errorf("invalid or missing cpu line in /proc/stat")
}*/


func GetCPUSample()(CPUSample, error){
	file, err := os.Open("/proc/stat")
	 
	if err != nil {
		return CPUSample{}, err //TODO add this logs of the tool 
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan(){
		fields := strings.Fields(scanner.Text())
		
		if len(fields) >= 5 && fields[metricIdentifier] == "cpu" { // len(fields) >= 5 make sure the program is not going to reach and out of bounds while parsing metrics
			var total uint64
			for _, val := range fields[userTime:]{
				n, _ := strconv.ParseUint(val, 10, 64)
				total += n // adds up all the time values in the
			}
			idle, _ := strconv.ParseUint(fields[idleTime], 10 , 64)
			
			return CPUSample{Idle: idle, Total: total}, nil
		}
	}
	return CPUSample{}, nil
}

func CalculateCPUusage()(uint64, error){
	
	sample, _ := GetCPUSample()
	total1 := sample.Total
	idle1 := sample.Idle

	time.Sleep(1 * time.Second)

	sample , _ = GetCPUSample()

	total2 := sample.Total
	idle2 := sample.Idle

	totalDelta := total2 - total1
	idleDelta :=  idle2 - idle1

	CPUusage := (1 - (idleDelta/totalDelta) ) * 100
	fmt.Print(CPUusage)
	return CPUusage, nil


	//return math.MaxUint64 ,nil
}

func GetRamUsage()(float64, error){
	
	var usedRamPercentage, total, available float64

	file, err := os.Open("/proc/meminfo")
	 
	if err != nil {
		return usedRamPercentage, err 
	}

	defer file.Close()
	
	scanner := bufio.NewScanner(file)
  	for scanner.Scan(){ //Scanner.Scan returns an double pointer array so it must be accesed as foo[x]
			line := scanner.Text()
			
			if strings.HasPrefix(line, "MemTotal"){

				fields := strings.Fields(line)
				total, _ = strconv.ParseFloat(fields[memoryParameter] , 64)
				
				
			}

			if strings.HasPrefix(line, "MemAvailable"){
				fields := strings.Fields(line)
				available, _ = strconv.ParseFloat(fields[memoryParameter], 64)

				
			}
			
			if total > 0 && available > 0 {
				break // break the for when total and available had been found
			}	

		}

		if err := scanner.Err(); err != nil {
    		log.Fatalf("Error encountered during file scanning: %v", err) //TODO change this from log to the tool log 
		}

		if total == 0 {
			fmt.Println(total, available, usedRamPercentage)
			return 0, nil
		}
	
		usedRamPercentage = ((total - available)/total) * 100
		fmt.Println(usedRamPercentage, total, available)
	

	return math.MaxUint64, nil // returns maxuint as the default value for an error
}