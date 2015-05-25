package output

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/fatih/color"
)

type Payload struct {
	Message string
	Error   bool
	Fatal   bool
}

func (p *Payload) IsError() {
	p.Error = true
}

func (p *Payload) IsFatal() {
	p.Fatal = true
}

func Error(message string) Payload {
	formatted := ApplyColour(message, "red")
	payload := Print(formatted)
	payload.IsError()
	return payload
}

func FatalError(message string) Payload {
	payload := Error(message)
	payload.IsFatal()
	return payload
}

func Print(message string) Payload {
	return Payload{message, false, false}
}

func ApplyColour(message string, outputType string) string {
	if runtime.GOOS != "windows" {
		switch outputType {
		case "red":
			return color.RedString(message)
		case "blue":
			return color.BlueString(message)
		}

	}
	return message
}

func Wait(channel chan Payload, done *sync.WaitGroup) {
	for {
		entry := <-channel
		fmt.Println(entry.Message)
		if entry.Fatal {
			done.Done()
			os.Exit(-1)
		}
	}
}
