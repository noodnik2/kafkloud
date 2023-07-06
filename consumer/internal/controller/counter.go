package controller

import (
	"log"
)

type Counter struct {
	nConsumed int
}

func (c *Counter) ProcessMessage(_, message string) error {
	log.Printf("processing message received(%s)\n", message)
	c.nConsumed += len(message)
	log.Printf("nConsumed(%d) after processing(%s)\n", c.nConsumed, message)
	return nil
}

func (c *Counter) GetNConsumed() int {
	return c.nConsumed
}
