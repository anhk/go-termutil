package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
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
	regexs = []*regexp.Regexp{
		// regexp.MustCompile("\x1b\\[[^@-~]*[@-~#]"),
		// regexp.MustCompile("\x1b[!-~#]*\x1b\\\\"),
		regexp.MustCompile(`\[[ -~]*][$#]*`),
		regexp.MustCompile(`\[[0-9#]*P`),
		regexp.MustCompile(`[ -~]*:[ -~]*#\s?`),
	}
)

type Processor struct {
	outTerm *termutil.Terminal
	inTerm  *termutil.Terminal

	output        []string
	outputAsInput bool
}

func (p *Processor) saveOutputLog() {
	fmt.Println("[OUTPUT]", strings.Join(p.output, "\n"))
	p.output = nil
}

func (p *Processor) processInputLog() {
	// fmt.Println("[INPUT]", p.inTerm.GetActiveBuffer().GetCurrentLine().String())

	if len(p.output) > 0 { // 如果output有暂存数据，则输出
		p.saveOutputLog()
	}
	p.outputAsInput = true
}

func (p *Processor) processOutputLog() {
	if p.outputAsInput {
		p.outputAsInput = false
		inputStr := p.outTerm.GetActiveBuffer().GetCurrentLine().String()
		for _, r := range regexs {
			inputStr = r.ReplaceAllString(inputStr, "")
		}
		if inputStr == "\n" || inputStr == "\r\n" {
			return
		}
		fmt.Println("[INPUT]", p.outTerm.GetActiveBuffer().GetCurrentLine().String())
		fmt.Println("[INPUT]", inputStr)
		// fmt.Printf("[INPUT]: %x", inputStr)
	} else {
		p.output = append(p.output, p.outTerm.GetActiveBuffer().GetCurrentLine().String())
		if len(p.output) > 1024 {
			p.saveOutputLog()
		}
	}
}

func main() {
	file, err := os.Open("./ssh.log")
	if err != nil {
		panic(err)
	}

	p := &Processor{}
	p.inTerm = termutil.New(termutil.WithRequestRender(p.processInputLog))
	p.outTerm = termutil.New(termutil.WithRequestRender(p.processOutputLog))

	r := bufio.NewReader(file)
	// terminal2.RequestRender = saveLog

	for {
		message, _, err := r.ReadLine()
		if err != nil {
			break
		}
		var lineResult Line
		json.Unmarshal([]byte(message), &lineResult)
		lineResult.TimeUnixNano = time.Now().UnixNano()
		b, _ := base64.StdEncoding.DecodeString(lineResult.Content)
		// fmt.Printf("=== %+v\n", lineResult)

		switch lineResult.Type {
		case "input":
			// terminal2.Write(b)
			// fmt.Printf("=== %+v\n", string(b))
			p.inTerm.Process(b)
		case "output":
			// terminal.Write(b)
			p.outTerm.Process(b)
		}
	}

	time.Sleep(time.Second * 2)
}
