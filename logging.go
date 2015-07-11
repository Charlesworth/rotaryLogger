package rotaryLogger

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

//StartLoggerWRotation starts the rotary logger. logAmount is the amount of
//log files that will be retained on a rotation. logLocation is the file and
//name (minus the .log prefix) of the logs, i.e. log/mylog. rotationPeriod
//is the time duration between log rotations
func StartLoggerWRotation(logAmount int, logLocation string, rotationPeriod time.Duration) {
	logRotator(logAmount, logLocation)

	StartLogger(logLocation)

	ticker := time.NewTicker(rotationPeriod)
	go func() {
		for _ = range ticker.C {
			Logger.Println("Rotating Logs")

			Logger = log.New(os.Stdout,
				"", //"PREFIX: ",
				log.Ldate|log.Ltime|log.Lshortfile)

			logRotator(logAmount, logLocation)

			StartLogger(logLocation)

		}
	}()
}

//StartLogger starts the Logger without any log rotation. logLocation is the
//file and folder location of the log, i.e. log/myAppLogs.
func StartLogger(logLocation string) {
	file, err := os.OpenFile(logLocation+"0.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open Logger file", err)
	}

	multi := io.MultiWriter(file, os.Stdout)

	Logger = log.New(multi,
		"", //"PREFIX: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

//logRotator is a function that deletes the latermost log and shifts all the
//rest of the logs by +1.
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
