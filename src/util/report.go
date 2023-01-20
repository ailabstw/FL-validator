package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var resultCorrect = true
var mux *sync.Mutex = new(sync.Mutex)

func MakeResultFalse() {
	mux.Lock()
	resultCorrect = false
	mux.Unlock()
}

func GetResult() bool {
	mux.Lock()
	r := resultCorrect
	mux.Unlock()
	return r
}

func WriteReport(state string, msg string, er string) {

	log.Println("writing report message ... state: " + state + ", msg: " + msg)
	err := os.MkdirAll(filepath.Dir(os.Getenv("REPORT_PATH")), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fo, err := os.OpenFile(os.Getenv("REPORT_PATH"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	var log ValidatingLogData
	log.Timestamp = time.Now().Format(time.RFC3339)
	log.State = state
	log.Message = msg
	log.Error = er
	jsonStr, _ := json.Marshal(log)
	fo.Write(jsonStr)
	fo.WriteString("\n")
	fo.Close()

}

type ValidatingLogData struct {
	Timestamp string `json:"timestamp"`
	State     string `json:"state"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}
