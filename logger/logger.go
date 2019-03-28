package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Out              *out
	logErrorPath     string
	logErrorFile     *os.File
	logErrorFilename string
)

type out struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

const (
	colorBlack = uint8(iota + 90)
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func Init(logPath string) {
	Out = &out{}

	flag := log.Ldate | log.Ltime | log.Lshortfile

	Out.Info = log.New(os.Stdout, Info(""), flag)
	Out.Warn = log.New(os.Stdout, Warn(""), flag)

	if logPath == "" {
		logPath = filepath.Dir(os.Args[0]) + "/log/"
	}

	logErrorPath = logPath

	if _, err := os.Stat(logErrorPath); err != nil {
		os.Mkdir(logErrorPath, os.ModePerm)
	}

	logErrorFilename = time.Now().Format("20060102")
	filename := logErrorPath + "error-" + logErrorFilename + ".txt"
	logErrorFile, _ = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	Out.Error = log.New(io.MultiWriter(logErrorFile, os.Stderr), Error(""), flag)

	checkDay()
}

func checkDay() {
	go func() {
		for range time.Tick(1 * time.Minute) {
			if filename := time.Now().Format("20060102"); filename != logErrorFilename {
				logErrorFile.Close()
				logErrorFilename = filename

				filename := logErrorPath + "error-" + logErrorFilename + ".txt"
				logErrorFile, _ = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
				Out.Error.SetOutput(io.MultiWriter(logErrorFile, os.Stderr))
			}
		}
	}()
}

func Trace(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorCyan, "TRACE", msg)
}

func Error(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorRed, "ERROR", msg)
}

func Warn(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorYellow, "WARN", msg)
}

func Info(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorGreen, "INFO", msg)
}

func Debug(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorBlue, "DEBUG", msg)
}

func Assert(msg string) string {
	return fmt.Sprintf("\x1b[%dm[%s] %s\x1b[0m", colorMagenta, "ASSERT", msg)
}
