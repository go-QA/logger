// Copyright 2013 The goQA Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.package goQA

package logger

import (
	"fmt"
	//"sync"
	//"error"
	"log"
	"runtime"
	//"os"
	"io"
	"time"
)

const (
	LOG_QUEUE_SIZE = 100
	LOG_SYNC_DELAY = 2
)

const LOG_LEVEL_NOT_SET = 0
const (
	LOG_LEVEL_DEBUG = (1 << iota)
	LOG_LEVEL_MESSAGE
	LOG_LEVEL_WARNING
	LOG_LEVEL_PASS
	LOG_LEVEL_FAIL
	LOG_LEVEL_RESULTS
	LOG_LEVEL_ERROR
	LOG_LEVEL_ALL
)

type logArg struct {
	level   uint64
	pattern string
	args    []interface{}
}

type logStream struct {
	ChnLogInput chan logArg
	level       uint64
	logger      log.Logger
}

func (log *logStream) Init(debug bool) {

	log.ChnLogInput = make(chan logArg, LOG_QUEUE_SIZE)

	go func(bool) {
		for message := range log.ChnLogInput {
			log.logger.Printf(message.pattern, message.args...)
		}
	}(debug)

}

func (log *logStream) sync() {
	for len(log.ChnLogInput) > 0 {
		time.Sleep(time.Millisecond * LOG_SYNC_DELAY)
	}
}

type GoQALog struct {
	chnInput    chan logArg
	loggers     map[string]logStream
	debugMode   bool
	initialized bool
	faulted     bool
	END         string
}

// Init method will automatically be called before logger is used but user can call if desired.
func (gLog *GoQALog) Init() {
	if gLog.initialized {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			gLog.faulted = true
			panic(r)
		}
	}()
	gLog.loggers = make(map[string]logStream)
	gLog.chnInput = make(chan logArg, LOG_QUEUE_SIZE)
	gLog.debugMode = false
	if runtime.GOOS == "linux" {
		gLog.END = "\n"
	} else {
		gLog.END = "\r\n"
	}
	gLog.initialized = true

	go func() {
		for message := range gLog.chnInput {
			for _, logger := range gLog.loggers {

				if ((message.level & LOG_LEVEL_DEBUG) != 0) &&
					((logger.level&(LOG_LEVEL_MESSAGE|LOG_LEVEL_ALL) != 0) && gLog.debugMode) {
					logger.ChnLogInput <- message
					continue
				}

				if (logger.level & LOG_LEVEL_ALL) != 0 {
					logger.ChnLogInput <- message
					continue
				}

				if ((logger.level & LOG_LEVEL_RESULTS) != 0) &&
					(message.level&(LOG_LEVEL_PASS|LOG_LEVEL_FAIL) != 0) {
					logger.ChnLogInput <- message
					continue
				}

				log_Level := uint64(message.level & logger.level)
				if log_Level != 0 {
					logger.ChnLogInput <- message
					continue
				}
			}
		}
	}()
}

func (gLog *GoQALog) ready() bool {
	if gLog.initialized {
		return true
	}
	if gLog.faulted {
		return false
	}
	gLog.Init()
	return gLog.initialized
}

func (gLog *GoQALog) Add(name string, level uint64, stream io.Writer) {
	if !gLog.ready() {
		return
	}
	if _, ok := gLog.loggers[name]; !ok {
		stream := logStream{level: level, logger: *log.New(stream, "", log.Ldate|log.Ltime|log.Lmicroseconds)}
		stream.Init(gLog.debugMode)
		gLog.loggers[name] = stream
	}
}

func (gLog *GoQALog) Printf(level uint64, value string, args ...interface{}) {
	arg := logArg{level, value, args}
	gLog.chnInput <- arg
}

func (gLog *GoQALog) Sync() {
	if !gLog.ready() {
		return
	}
	for len(gLog.chnInput) > 0 {
		time.Sleep(time.Millisecond * LOG_SYNC_DELAY)
	}
	for _, logger := range gLog.loggers {
		logger.sync()
	}
}
func (gLog *GoQALog) SetDebug(mode bool) {
	if gLog.ready() {
		gLog.debugMode = mode
	}
}

func (gLog *GoQALog) IsDebugSet() bool {

	return gLog.debugMode
}

func (gLog *GoQALog) LogError(errMsg string, args ...interface{}) {
	if gLog.ready() {
		gLog.Printf(LOG_LEVEL_ERROR, fmt.Sprintf("ERROR::%s%s", errMsg, gLog.END), args...)
	}
}

func (gLog *GoQALog) LogDebug(DebugMsg string, args ...interface{}) {
	if gLog.ready() {
		if gLog.debugMode {
			gLog.Printf(LOG_LEVEL_DEBUG, fmt.Sprintf("DEBUG::%s%s", DebugMsg, gLog.END), args...)
		}
	}
}

func (gLog *GoQALog) LogWarning(warnMsg string, args ...interface{}) {
	if gLog.ready() {
		gLog.Printf(LOG_LEVEL_WARNING, fmt.Sprintf("ERROR::%s%s", warnMsg, gLog.END), args...)
	}
}

func (gLog *GoQALog) LogPass(passMsg string, args ...interface{}) {
	if gLog.ready() {
		gLog.Printf(LOG_LEVEL_PASS, fmt.Sprintf("PASS::%s%s", passMsg, gLog.END), args...)
	}
}

func (gLog *GoQALog) LogFail(failMsg string, args ...interface{}) {
	if gLog.ready() {
		gLog.Printf(LOG_LEVEL_FAIL, fmt.Sprintf("FAIL::%s%s", failMsg, gLog.END), args...)
	}
}

func (gLog *GoQALog) LogMessage(msg string, args ...interface{}) {
	if gLog.ready() {
		gLog.Printf(LOG_LEVEL_MESSAGE, fmt.Sprintf("MSG::%s%s", msg, gLog.END), args...)
	}
}
