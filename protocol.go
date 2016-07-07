// Copyright 2016 The Godaq Authors. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package godaq

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type CommandNumber uint8

const nak = 160

var (
	ErrChecksum      = errors.New("Checksum error")
	ErrInvalidLength = errors.New("Invalid message length")
	ErrNakReceived   = errors.New("NAK response received")
)

type Message struct {
	Number CommandNumber
	Body   []byte
}

func checksum(data []byte) (csum uint16) {
	for _, b := range data {
		csum += uint16(b)
	}
	return
}

func toBytes(value interface{}) []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, value)
	return b.Bytes()
}

func (m *Message) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(m.Body))
	b[2] = byte(m.Number)
	b[3] = byte(len(m.Body))
	copy(b[4:], m.Body)

	// place the checksum at the start of the message
	binary.BigEndian.PutUint16(b[:2], checksum(b))
	return b, nil
}

func parseResponse(b []byte) (io.Reader, error) {
	csum := checksum(b[2:])
	if binary.BigEndian.Uint16(b[:2]) != csum {
		return nil, ErrChecksum
	}
	if b[2] == nak {
		return nil, ErrNakReceived
	}
	if int(b[3]) != len(b)-4 {
		return nil, ErrInvalidLength
	}
	return bytes.NewBuffer(b[4:]), nil
}

func sendCommand(ser io.ReadWriter, command *Message, respLen int) (io.Reader, error) {
	data, err := command.Marshal()
	if err != nil {
		return nil, err
	}
	if _, err := ser.Write(data); err != nil {
		return nil, err
	}
	data = make([]byte, respLen+4)
	if _, err := ser.Read(data); err != nil {
		return nil, err
	}
	return parseResponse(data)
}
