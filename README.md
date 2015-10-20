# logger
Logger object used in goQA. Wrapper around Go logger that handles multiple log entities and log levels. Thread safe

##Features:
- Handle multiple log streams. Any object using io.Write interface
 
## Quick Start

To download, run the following command:

~~~
go get github.com/go-QA/logger
~~~


Open some files for writing logs:
```go
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
```

Create a log object and call the `Init()` method:

```go
	log := logger.GoQALog{}
	log.Init()
	log.SetDebug(true)
```

  Add logs to log object using `Add(name string, level uint64, stream io.Writer)` Now we have for log files and output to `os.Stdout`
```go
	log.Add("default", logger.LOGLEVEL_ALL, os.Stdout)
	log.Add("Console", logger.LOGLEVEL_MESSAGE, console)
	log.Add("Error", logger.LOGLEVEL_ERROR, errLog)
	log.Add("Incidents", logger.LOGLEVEL_WARNING, incedentLog)
	log.Add("Resuts", logger.LOGLEVEL_RESULTS, resultLog)
```

  Use APIs to appropriate outputs depending on the log level chosen. Methods defined same as Printf, `LogMessage(msg string, args ...interface{})`

```go
	log.LogMessage("running on platform %s", runtime.GOOS)
	log.LogMessage("First message")
	log.LogMessage("second message")
	log.LogMessage("third message")
	log.LogDebug("Debug message")
	log.LogWarning("Warning Will Robinson")
	log.LogPass("Test Passed")
	log.LogFail("Test Failed")
	log.LogError("Failure in script")
```




