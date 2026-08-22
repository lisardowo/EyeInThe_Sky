package helpers

import (
	"os"
)

var ReportFile *os.File

func InitReport(path string) error {
	fd, err := os.OpenFile(path, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	ReportFile = fd
	
	return nil
}

func CloseReport(){
	if ReportFile != nil{
		ReportFile.Close()
	}
}

func WriteReport(entry LogEntry ){
	if ReportFile == nil {
		return
	}
	ReportFile.WriteString(entry.FormatEntry() + "\n")
}