package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	flag        = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	globalLevel = LEVEL_TRACE
	logFiles    map[int]*os.File
	logWrites   map[int]*log.Logger
	levels      map[int]string
	config      *Config
)

// TRACE < DEBUG < INFO < WARN < ERROR
const (
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

const (
	colorTrace = "1;32"
	colorDebug = "1;32"
	colorInfo  = "1;36"
	colorWarn  = "1;33"
	colorError = "1;31"
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

	logFiles = make(map[int]*os.File)
	logWrites = make(map[int]*log.Logger)

	initLoggerWrite()

	checkDay()
}

// 每天生成一个文件
func checkDay() {
	go func() {
		for range time.Tick(1 * time.Minute) {
			initLoggerWrite()
		}
	}()
}

func initLoggerWrite() {
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

			logWrites[i] = log.New(io.MultiWriter(logFiles[i], os.Stdout), "", flag)
		} else {
			logWrites[i] = log.New(os.Stdout, "", flag)
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

func Trace(msg string) {
	if LEVEL_TRACE >= globalLevel {
		logWrites[LEVEL_TRACE].Println(fmt.Sprintf("\033[%sm[%s] %s\033[0m", colorTrace, "TRACE", msg))
	}
}

func Debug(msg string) {
	if LEVEL_DEBUG >= globalLevel {
		logWrites[LEVEL_DEBUG].Println(fmt.Sprintf("\033[%sm[%s] %s\033[0m", colorDebug, "DEBUG", msg))
	}
}

func Info(msg string) {
	if LEVEL_INFO >= globalLevel {
		logWrites[LEVEL_INFO].Println(fmt.Sprintf("\033[%sm[%s] %s\033[0m", colorInfo, "INFO", msg))
	}
}

func Warn(msg string) {
	if LEVEL_WARN >= globalLevel {
		logWrites[LEVEL_WARN].Println(fmt.Sprintf("\033[%sm[%s] %s\033[0m", colorWarn, "WARN", msg))
	}
}

func Error(msg string) {
	if LEVEL_ERROR >= globalLevel {
		logWrites[LEVEL_ERROR].Println(fmt.Sprintf("\033[%sm[%s] %s\033[0m", colorError, "ERROR", msg))
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
