package main

import (
	"io"
	"os"
	"schedule-timer/internal/model"
	"schedule-timer/internal/reader/yaml"
	"schedule-timer/internal/ui/tui"
)

func main() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(1)
	}

	project, err := readProjectFromFile(os.Args[1])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	err = tui.Run(*project)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func usage() {
	println("Usage: schedule-timer <file>")
}

func readProjectFromFile(filePath string) (*model.Project, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return yaml.Read(data)
}
