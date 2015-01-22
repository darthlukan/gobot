package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func ChannelLogger(Log string, UserNick string, message string) {
	STime := time.Now().UTC().Format(time.ANSIC)
	log := strings.Replace(Log, "#", "", 1)
	logFile := fmt.Sprintf("%s.log", log)

	//Open the file for writing With Append Flag to create file persistence
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	n, err := io.WriteString(f, fmt.Sprintf("%v > %v: %v\n", STime, UserNick, message))
	if err != nil {
		fmt.Println(n, err)
	}
}

func LogDir(CreateDir string) {
	if _, err := os.Stat(CreateDir); os.IsNotExist(err) {
		fmt.Printf("No such file or directory: %s\n", CreateDir)
		os.Mkdir(CreateDir, 0777)
	} else {
		fmt.Printf("Its There: %s\n", CreateDir)
	}
}

func LogFile(CreateFile string) {
	log := strings.Replace(CreateFile, "#", "", 1)
	logFile := fmt.Sprintf("%s.log", log)

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Printf("Log File %s Doesn't Exist. Creating Log File.\n", logFile)
		os.Create(logFile)
		fmt.Printf("Log File %s Created.\n", logFile)
	} else {
		fmt.Printf("Log File %s Exists.\n", logFile)
	}
}
