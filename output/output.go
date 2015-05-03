package output

import (
	"fmt"
	"os"
	"runtime"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/fatih/color"
)

type Payload struct {
	Message string
	Error   bool
}

func (p *Payload) IsError() {
	p.Error = true
}

func Error(message string) Payload {
	formatted := ApplyColour(message, "red")
	payload := Print(formatted)
	payload.IsError()
	return payload
}

func Print(message string) Payload {
	return Payload{message, false}
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

func Wait(channel chan Payload) {
	for {
		entry := <-channel
		if entry.Error {
			os.Exit(-1)
		} else {
			fmt.Println(entry.Message)
		}
	}
}
