package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/anhk/go-termutil"
)

type Line struct {
	Type         string `json:"type"`
	ClientIp     string `json:"client_ip"`
	Ip           string `json:"ip"`
	User         string `json:"user"`
	Session      string `json:"session"`
	Content      string `json:"content"`
	Timestamp    int64  `json:"timestamp"`
	TimeUnixNano int64  `json:"time_unix_nano"`
}

var (
	terminal = termutil.New()
	// terminal2  = termutil.New()
	lineResult Line
	// outputAsInput bool
)

func saveLog() {
	// if outputAsInput {
	// 	lineResult.Type = "input"
	// 	outputAsInput = false
	// }
	fmt.Println(terminal.GetActiveBuffer().GetCurrentLine().String())
	// fmt.Println(terminal2.GetActiveBuffer().GetCurrentLine().String())
}

func main() {
	file, err := os.Open("./ssh.log")
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)
	terminal.RequestRender = saveLog
	// terminal2.RequestRender = saveLog

	for {
		message, _, err := r.ReadLine()
		if err != nil {
			break
		}
		json.Unmarshal([]byte(message), &lineResult)
		lineResult.TimeUnixNano = time.Now().UnixNano()
		b, _ := base64.StdEncoding.DecodeString(lineResult.Content)
		// fmt.Printf("=== %+v\n", lineResult)
		switch lineResult.Type {
		case "input":
			// terminal2.Write(b)
			// fmt.Printf("=== %+v\n", string(b))

		case "output":
			// terminal.Write(b)
			terminal.Process(b)
		}
	}

	time.Sleep(time.Second * 2)
}
