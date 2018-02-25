package envsensor

import (
	"fmt"
	"testing"
)

func TestNewReadingBroadcast(t *testing.T) {
	input := make(chan Reading)
	b := NewReadingBroadcast(input)

	if b.Input != input {
		t.Fatal("No input set")
	}
}

func TestReadingBroadcastStartNoOutput(t *testing.T) {
	input := make(chan Reading)
	b := NewReadingBroadcast(input)

	if len(b.Outputs) != 0 {
		t.Fatal(fmt.Sprintf("Found %d outputs, expected 0", len(b.Outputs)))
	}

	err := b.Start()
	if err == nil {
		t.Fatal("No error received for starting with no outputs")
	}
}

func TestReadingBroadcastAddOutput(t *testing.T) {
	input := make(chan Reading)
	b := NewReadingBroadcast(input)
	if len(b.Outputs) != 0 {
		t.Fatal(fmt.Sprintf("Found %d outputs, expected 0", len(b.Outputs)))
	}

	output := make(chan Reading)
	b.AddOutput(output)

	if len(b.Outputs) != 1 {
		t.Fatal(fmt.Sprintf("Found %d outputs, expected 1", len(b.Outputs)))
	}
}

func TestReadingBroadcastStartStop(t *testing.T) {
	input := make(chan Reading, 2)
	b := NewReadingBroadcast(input)
	output := make(chan Reading)
	b.AddOutput(output)

	// fakeReading := Reading{}

	var err error
	err = b.Start()
	if err != nil {
		t.Fatal("Got error starting broadcast", err)
	}
	err = b.Stop()
	if err != nil {
		t.Fatal("Got error stopping broadcast", err)
	}

	// input <- fakeReading

}
