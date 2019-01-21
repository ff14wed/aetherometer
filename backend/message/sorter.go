package message

import (
	"bufio"
	"bytes"
)

type mailbox struct {
	buffer *bytes.Buffer
	reader *bufio.Reader
}

func newMailbox() mailbox {
	buffer := bytes.NewBuffer(nil)
	return mailbox{
		buffer: buffer,
		reader: bufio.NewReaderSize(buffer, 65536),
	}
}

// Sorter sorts messages to the correct destination.
// This is especially useful for handling multiple sources of data and feeding
// them to the correct handler in a single goroutine. This same task could be
// done with goroutines, but it may lead to goroutine leaks when it is not clear
// or impossible to know when these sources of data appear or disappear.
type Sorter struct {
	mailboxes map[uint32]mailbox
}

// NewSorter creates a new sorter
func NewSorter() *Sorter {
	return &Sorter{
		mailboxes: make(map[uint32]mailbox),
	}
}

// Sort sorts each message into the correct buffer and returns the buffered
// reader that received the new message.
// If the buffer doesn't exist, it creates a new one.
func (s *Sorter) Sort(dst uint32, message []byte) *bufio.Reader {
	if _, found := s.mailboxes[dst]; !found {
		s.mailboxes[dst] = newMailbox()
	}
	_, _ = s.mailboxes[dst].buffer.Write(message)
	return s.mailboxes[dst].reader
}
