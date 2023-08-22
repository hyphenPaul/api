package main

import (
	"encoding/json"
	"fmt"
)

type loggerPayload struct {
	Duration string `json:"duration"`
	URL      string `json:"url"`
	Method   string `json:"method"`
}

type logger interface {
	info(l loggerPayload)
	error(error)
}

type noopLogger struct{}

func (j noopLogger) info(p loggerPayload) {}
func (j noopLogger) error(e error)        {}

type jsonLogger struct{}

func (j jsonLogger) info(p loggerPayload) {
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(fmt.Errorf("logger error: %v\n", err))
		return
	}

	fmt.Println(string(b))
}

func (j jsonLogger) error(e error) {
	fmt.Println("{\"error\" : \"", e.Error(), "\"}")
}
