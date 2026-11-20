package protocol

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

const MaxMessageSize = 64 * 1024

type Decoder struct {
	scanner *bufio.Scanner
}

func NewDecoder(reader io.Reader) *Decoder {
	scanner := bufio.NewScanner(reader)

	scanner.Buffer(
		make([]byte, 4096),
		MaxMessageSize,
	)

	return &Decoder{
		scanner: scanner,
	}
}

func (d *Decoder) Decode() (ClientMessage, error) {
	if !d.scanner.Scan() {
		if err := d.scanner.Err(); err != nil {
			return ClientMessage{}, fmt.Errorf("read message: %w", err)
		}

		return ClientMessage{}, io.EOF
	}

	var message ClientMessage

	if err := json.Unmarshal(d.scanner.Bytes(), &message); err != nil {
		return ClientMessage{}, fmt.Errorf("decode message: %w", err)
	}

	return message, nil
}

type Encoder struct {
	encoder *json.Encoder
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		encoder: json.NewEncoder(writer),
	}
}

func (e *Encoder) Encode(message ServerMessage) error {
	if err := e.encoder.Encode(message); err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	return nil
}


type ClientEncoder struct {
	encoder *json.Encoder
}

func NewClientEncoder(writer io.Writer) *ClientEncoder {
	return &ClientEncoder{
		encoder: json.NewEncoder(writer),
	}
}

func (e *ClientEncoder) Encode(message ClientMessage) error {
	return e.encoder.Encode(message)
}


type ServerDecoder struct {
	scanner *bufio.Scanner
}

func NewServerDecoder(reader io.Reader) *ServerDecoder {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), MaxMessageSize)

	return &ServerDecoder{
		scanner: scanner,
	}
}

func (d *ServerDecoder) Decode() (ServerMessage, error) {
	var message ServerMessage

	if !d.scanner.Scan() {
		if err := d.scanner.Err(); err != nil {
			return message, err
		}

		return message, io.EOF
	}

	if err := json.Unmarshal(d.scanner.Bytes(), &message); err != nil {
		return message, err
	}

	return message, nil
}
