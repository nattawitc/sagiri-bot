package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "sagiri-bot/config"

	"github.com/spf13/viper"
)

const ()

var (
	errorPath    string
	errorChannel chan string
	LogWaitGroup sync.WaitGroup
)

func init() {
	errorPath = viper.GetString("log-path.error")
	if errorPath == "" {
		panic("log path is missing")
	}
	errorChannel = make(chan string)
	go startLogger(errorChannel, errorPath)
}

func PrintError(a ...interface{}) {
	_, path, lineNumber, _ := runtime.Caller(1)
	paths := strings.Split(path, "/")
	filename := fmt.Sprintf("%v(%v) ", paths[len(paths)-1], lineNumber)
	LogWaitGroup.Add(1)
	errorChannel <- filename + fmt.Sprint(a...)
}

func startLogger(ch chan string, logPath string) {
	err := createDirIfNotExist(logPath)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		func() {
			logMsg := <-ch
			defer LogWaitGroup.Done()
			currentTime := time.Now()
			filename := fmt.Sprintf("%v/%04d-%02d.log", logPath, currentTime.Year(), currentTime.Month())
			outfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				log.Println("can't open file", filename)
				log.Println(logMsg)
				return
			}
			defer outfile.Close()
			mul := io.MultiWriter(os.Stdout, outfile)
			logger := log.New(mul, "", log.LstdFlags)
			logger.Println(logMsg)
		}()
	}
}

func createDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func dumpRequestBody(req *http.Request) ([]byte, error) {
	body, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	lastIndex := bytes.LastIndex(body, []byte("\r\n\r\n")) + 4
	return body[lastIndex:], nil
}
