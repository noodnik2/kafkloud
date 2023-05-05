package controller

import (
	"log"
)

type Controller interface {
	ProcessMessage(input string)
	GetNConsumed() int // number of words processed
}

type SimpleController struct {
	nConsumed int
}

func NewController() Controller {
	return &SimpleController{}
}

func (c *SimpleController) ProcessMessage(message string) {
	log.Printf("processing message received(%s)\n", message)
	c.nConsumed += len(message)
	log.Printf("nConsumed(%d) after processing(%s)\n", c.nConsumed, message)
}

func (c *SimpleController) GetNConsumed() int {
	return c.nConsumed
}
