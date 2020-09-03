package fileutil

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/erock530/go.logging"
)

// PipeStruct specifies a slice of bytes; used to represent piped stdin data.
type PipeStruct struct {
	Data []byte `json:"data"`
}

// ReadFromStdin reads from stdin until requested to stop.
func ReadFromStdin(stop <-chan struct{}) <-chan []byte {
	outChan := make(chan []byte)
	go func() {
		jReader := json.NewDecoder(os.Stdin)
		p := PipeStruct{}
		var err error
		for err == nil {
			err = jReader.Decode(&p)
			select {
			case outChan <- p.Data:
			case <-stop:
				return
			}
		}
		if err != nil {
			logging.Errorf("Error reading from stdin: %+v", err)
		}
		close(outChan)
	}()
	return outChan
}

// StdinReader is a collection consisting of a JSON decoder and slice of bytes.
type StdinReader struct {
	jReader *json.Decoder
	buf     []byte
}

// NewStdinReader creates and initializes a StdinReader struct.
func NewStdinReader() StdinReader {
	s := StdinReader{}
	s.buf = make([]byte, 0)
	s.jReader = json.NewDecoder(os.Stdin)
	return s
}

// Read decodes incoming JSON data and returns the number of bytes read.
func (s StdinReader) Read(p []byte) (n int, err error) {
	for len(s.buf) < len(p) && err == nil {
		// Grab the next packet and append the bytes to our buffer
		ps := PipeStruct{}
		err = s.jReader.Decode(&ps)
		s.buf = append(s.buf, ps.Data...)
	}

	n = copy(p, s.buf)
	s.buf = s.buf[n:]
	return
}

// WriteToStdout encodes incoming data and pipes it to stdout.
func WriteToStdout(outChan <-chan []byte) error {
	p := PipeStruct{}
	jWriter := json.NewEncoder(os.Stdout)
	var d []byte
	var i int
	for d = range outChan {
		p.Data = d
		err := jWriter.Encode(p)
		if err != nil {
			return fmt.Errorf("Error writing to stdout: %+v", err)
		}
		i++
	}
	return nil
}

// StdoutWriter represents an JSON encoder.
type StdoutWriter struct {
	jWriter *json.Encoder
}

// NewStdoutWriter creates and initializes a StdoutWriter struct.
func NewStdoutWriter() StdoutWriter {
	return StdoutWriter{jWriter: json.NewEncoder(os.Stdout)}
}

// Write encodes outgoing JSON data and returns the number of bytes written.
func (s StdoutWriter) Write(p []byte) (int, error) {
	n := len(p)
	return n, s.jWriter.Encode(PipeStruct{Data: p})
}
