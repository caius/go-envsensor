package envsensor

import (
	"errors"
)

type ReadingBroadcast struct {
	Input    chan Reading
	Outputs  []chan Reading
	shutdown chan string
}

func NewReadingBroadcast(input chan Reading) ReadingBroadcast {
	return ReadingBroadcast{Input: input}
}

func (rb *ReadingBroadcast) AddOutput(output chan Reading) {
	rb.Outputs = append(rb.Outputs, output)
}

func (rb *ReadingBroadcast) Start() error {
	if len(rb.Outputs) == 0 {
		return errors.New("Outputs required before starting broadcast")
	}

	return nil
}
