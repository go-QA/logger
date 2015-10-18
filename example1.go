package main

import (
	//"fmt"
	//"error"
	//"log"
	"os"
	"runtime"
	//"io"
	"time"
	//"net"
	//"encoding/json"
	"github.com/go-QA/logger"
)

func main() {
	
 	console, err := os.Create("data/console.log")
	if err != nil { panic(err) }
	defer console.Close()

	errLog, err := os.Create("data/error.log")
	if err != nil { panic(err) }
	defer errLog.Close()

	incedentLog, err := os.Create("data/incedents.log")
	if err != nil { panic(err) }
	defer incedentLog.Close()

	resultLog, err := os.Create("data/TestResults.log")
	if err != nil { panic(err) }
	defer resultLog.Close()

	log := goQA.GoQALog{}
	log.Init()
	log.SetDebug(true)

	log.Add("default", logger.LOGLEVEL_ALL, os.Stdout)
	log.Add("Console", logger.LOGLEVEL_MESSAGE, console)
	log.Add("Error", logger.LOGLEVEL_ERROR, errLog)
	log.Add("Incidents", logger.LOGLEVEL_WARNING, incedentLog)
	log.Add("Resuts", logger.LOGLEVEL_PASS_FAIL, resultLog)

	log.LogMessage("running on platform %s", runtime.GOOS)
	log.LogMessage("First message")
	log.LogMessage("second message")
	log.LogMessage("third message")
	log.LogDebug("Debug message")
	log.LogWarning("Warning Will Robinson")
	log.LogPass("Test Passed")
	log.LogFail("Test Failed")
	log.LogError("Failure in script")

	time.Sleep(time.Second * 1)
}
