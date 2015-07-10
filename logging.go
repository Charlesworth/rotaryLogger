package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

//Logger is the logging object that is used to produce log outputs
var Logger *log.Logger

var file os.File

func startLogger(logAmount int, logLocation string, rotationPeriod time.Duration) {
	logRotator(logAmount, logLocation)

	newLogFile(logLocation)

	ticker := time.NewTicker(rotationPeriod)
	go func() {
		for _ = range ticker.C {
			Logger.Println("poo")

			Logger = log.New(os.Stdout,
				"", //"PREFIX: ",
				log.Ldate|log.Ltime|log.Lshortfile)

			logRotator(logAmount, logLocation)

			newLogFile(logLocation)
			Logger.Println("wee")

		}
	}()
}

func newLogFile(logLocation string) {
	file, err := os.OpenFile(logLocation+"0.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open Logger file", err)
	}

	multi := io.MultiWriter(file, os.Stdout)

	Logger = log.New(multi,
		"", //"PREFIX: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func logRotator(logAmount int, logLocation string) {
	//if log file number <logAmount> exists, delete it
	if _, err := os.Stat(logLocation + strconv.Itoa(logAmount) + ".log"); os.IsNotExist(err) {
		os.Remove(logLocation + strconv.Itoa(logAmount) + ".log")
	}

	//shift all of the other logs along by +1
	for i := 2; i != -1; i-- {
		fmt.Println(logLocation + strconv.Itoa(i) + ".log")
		if _, err := os.Stat(logLocation + strconv.Itoa(i) + ".log"); err == nil {
			os.Rename(logLocation+strconv.Itoa(i)+".log", logLocation+strconv.Itoa(i+1)+".log")
		}
	}
}
