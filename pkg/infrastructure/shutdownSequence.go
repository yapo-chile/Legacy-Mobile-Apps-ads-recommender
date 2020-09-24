package infrastructure

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
)

// ShutdownSequence is a stack implementation to control the shutdown order of each
// infrastructure component needed. First process that starts is the last to be turned off
type ShutdownSequence struct {
	sequence  []io.Closer
	waitGroup *sync.WaitGroup
}

// Push pushes a new component into the stack to be turned off.
func (s *ShutdownSequence) Push(task io.Closer) {
	s.sequence = append([]io.Closer{task}, s.sequence...)
	s.waitGroup.Add(1)
}

// NewShutdownSequence creates a new ShutdownSequence
func NewShutdownSequence() *ShutdownSequence {
	var sequence []io.Closer
	var waitGroup sync.WaitGroup
	return &ShutdownSequence{
		sequence:  sequence,
		waitGroup: &waitGroup,
	}
}

// Wait waits until the internal waitGroup counter is zero.
func (s *ShutdownSequence) Wait() {
	s.waitGroup.Wait()
}

// close does the actual closing of things
func (s *ShutdownSequence) close() {
	for i := range s.sequence {
		if err := s.sequence[i].Close(); err != nil {
			fmt.Printf("Error closing the task of type %T: %+v\n", s.sequence[i].Close(), err)
		}
		s.waitGroup.Done()
	}
}

// Listen launches a go routines that waits for sigint and then stops each task in the stack.
// You need to call Listen before calling Wait, otherwise you risk waiting indefinitely
func (s *ShutdownSequence) Listen() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint)
		loop := true
		for loop {
			sig := <-sigint
			fmt.Printf("Received signal: %s\n", sig)
			if sig == os.Interrupt {
				// We received an interrupt signal, shut down.
				fmt.Printf("Proceeding to shut down\n")
				loop = false
			}
		}
		s.close()
		// At this point all processes must be done
	}()
}
