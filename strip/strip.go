package strip

import (
	"errors"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"sync"
)

var Status = map[int]string{
	0: "ok",
	1: "timed out",
	2: "invalid command",
	3: "invalid mode",
	4: "internal error",
}

const (
	StatusOK             = 0
	StatusTimeout        = 1
	StatusInvalidCommand = 2
	StatusInvalidMode    = 3
	StatusInternalError  = 4
)

var ErrSizeMismatch = errors.New("size mismatch")
var ErrWrite = errors.New("write error")

type Strip struct {
	Options serial.OpenOptions

	port      io.ReadWriteCloser
	portMutex sync.Mutex
}

type stripResponse struct {
	buffer []byte
	err    error
}

const responseSize = 5

func (s *Strip) Initialize() error {
	var err error
	s.port, err = serial.Open(s.Options)
	return err
}

func (s *Strip) getResponseChan() chan *stripResponse {
	responseChan := make(chan *stripResponse)

	go func() {
		buffer := make([]byte, responseSize)
		read, err := s.port.Read(buffer)

		if err == nil && read < responseSize {
			err = ErrSizeMismatch
		}

		responseChan <- &stripResponse{
			buffer: buffer,
			err:    err,
		}
	}()

	return responseChan
}

func (s *Strip) Execute(command []byte) error {
	s.portMutex.Lock()
	defer s.portMutex.Unlock()

	written, err := s.port.Write(command)
	if err != nil {
		return err
	}

	if written != len(command) {
		return ErrWrite
	}

	return nil
}

func (s *Strip) Request(command []byte) ([]byte, error) {
	s.portMutex.Lock()
	defer s.portMutex.Unlock()

	responseChan := s.getResponseChan()

	written, err := s.port.Write(command)
	if err != nil {
		return nil, err
	}

	if written != len(command) {
		return nil, ErrWrite
	}

	response := <-responseChan

	return response.buffer, response.err
}
