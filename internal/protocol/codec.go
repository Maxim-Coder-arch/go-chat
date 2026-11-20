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
    writer io.Writer
}

func NewEncoder(writer io.Writer) *Encoder {
    return &Encoder{
        writer: writer,
    }
}

func (e *Encoder) Encode(message ServerMessage) error {
    encoder := json.NewEncoder(e.writer)

    if err := encoder.Encode(message); err != nil {
        return fmt.Errorf("encode message: %w", err)
    }

    return nil
}