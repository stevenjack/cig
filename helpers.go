package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/stevenjack/cig/Godeps/_workspace/src/github.com/fatih/color"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func error_output(message string) string {
	return print_output(message, "red")
}

func output(channel chan string) {
	for {
		entry := <-channel
		fmt.Printf(entry)
	}
}

func print_output(message string, output_type string) string {
	if runtime.GOOS != "windows" {
		switch output_type {
		case "red":
			return color.RedString(message)
		case "blue":
			return color.BlueString(message)
		}

	}
	return message
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
