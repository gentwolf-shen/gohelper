package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	flag        = log.Ldate | log.Ltime | log.Lmicroseconds
	globalLevel = LEVEL_TRACE
	config      *Config
	logFiles    map[int]*os.File
	logWriters  map[int]*log.Logger
	levels      map[int]string
	colors      map[int]string
)

// TRACE < DEBUG < INFO < WARN < ERROR
const (
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

// 从配置文件中初始化
func InitFromJson(filename string) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("init logger config file error " + err.Error())
	}

	initLogger(b)
}

// 使用默认初始化
func InitDefault() {
	str := `{
			  "level": "DEBUG",
			  "file": {
				"logPath": "./log/",
				"level": "WARN"
			  }
			}`

	initLogger([]byte(str))
}

func initLogger(bytes []byte) {
	config = &Config{}
	if err := json.Unmarshal(bytes, config); err != nil {
		panic("parse logger config error " + err.Error())
	}

	if _, err := os.Stat(config.File.LogPath); err != nil {
		if err := os.MkdirAll(config.File.LogPath, os.ModePerm); err != nil {
			panic("create file.logPath error, " + err.Error())
		}
	}

	levels = make(map[int]string)
	levels[LEVEL_TRACE] = "TRACE"
	levels[LEVEL_DEBUG] = "DEBUG"
	levels[LEVEL_INFO] = "INFO"
	levels[LEVEL_WARN] = "WARN"
	levels[LEVEL_ERROR] = "ERROR"

	colors = make(map[int]string)
	colors[LEVEL_TRACE] = "1;32"
	colors[LEVEL_DEBUG] = "1;32"
	colors[LEVEL_INFO] = "1;36"
	colors[LEVEL_WARN] = "1;33"
	colors[LEVEL_ERROR] = "1;31"

	logFiles = make(map[int]*os.File)
	logWriters = make(map[int]*log.Logger)

	initLoggerWriter()

	checkDay()
}

// 每天生成一个文件
func checkDay() {
	go func() {
		for range time.Tick(1 * time.Minute) {
			initLoggerWriter()
		}
	}()
}

func initLoggerWriter() {
	// 关闭所有已打开的文件
	for _, file := range logFiles {
		file.Close()
	}

	// 输出到文件的最低级别
	minFileLevel := getMinFileLevel()

	// 创建日志输出的文件
	var err error
	date := time.Now().Format("20060102")
	for i, level := range levels {
		if i >= minFileLevel {
			filename := fmt.Sprintf("%s/%s_%s.log", config.File.LogPath, strings.ToLower(level), date)
			logFiles[i], err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

			if err != nil {
				panic("create log file error: " + filename)
			}

			logWriters[i] = log.New(io.MultiWriter(logFiles[i], os.Stdout), "", flag)
		} else {
			logWriters[i] = log.New(os.Stdout, "", flag)
		}
	}
}

func getMinFileLevel() int {
	minFileLevel := LEVEL_WARN
	for i, level := range levels {
		if level == config.File.Level {
			minFileLevel = i
			break
		}
	}

	return minFileLevel
}

// 动态设置日志输出级别
func SetLevel(level int) {
	if level < LEVEL_TRACE || level > LEVEL_ERROR {
		level = LEVEL_DEBUG
	}
	globalLevel = level
}

func GetLevel() int {
	return globalLevel
}

func Trace(msg interface{}) {
	logWriter(LEVEL_TRACE, msg)
}

func Tracef(format string, msg ...interface{}) {
	logWriter(LEVEL_TRACE, fmt.Sprintf(format, msg...))
}

func Debug(msg interface{}) {
	logWriter(LEVEL_DEBUG, msg)
}

func Debugf(format string, msg ...interface{}) {
	logWriter(LEVEL_DEBUG, fmt.Sprintf(format, msg...))
}

func Info(msg interface{}) {
	logWriter(LEVEL_INFO, msg)
}

func Infof(format string, msg ...interface{}) {
	logWriter(LEVEL_INFO, fmt.Sprintf(format, msg...))
}

func Warn(msg interface{}) {
	logWriter(LEVEL_WARN, msg)
}

func Warnf(format string, msg ...interface{}) {
	logWriter(LEVEL_WARN, fmt.Sprintf(format, msg...))
}

func Error(msg interface{}) {
	logWriter(LEVEL_ERROR, msg)
}

func Errorf(format string, msg ...interface{}) {
	logWriter(LEVEL_ERROR, fmt.Sprintf(format, msg...))
}

func Log(level int, msg interface{}) {
	logWriter(level, msg)
}
func Logf(level int, format string, msg ...interface{}) {
	logWriter(level, fmt.Sprintf(format, msg...))
}

func logWriter(level int, msg interface{}) {
	if level > globalLevel {
		if _, file, line, ok := runtime.Caller(2); ok {
			segments := strings.Split(file, "/")
			filename := segments[len(segments)-1]

			if level >= globalLevel {
				logWriters[level].Println(fmt.Sprintf("\033[%sm%s:%d [%s] %v\033[0m", colors[level], filename, line, levels[level], msg))
			}
		}
	}
}

type (
	Config struct {
		Level string `json:"level"`
		File  struct {
			LogPath string `json:"logPath"`
			Level   string `json:"level"`
		} `json:"file"`
	}
)
