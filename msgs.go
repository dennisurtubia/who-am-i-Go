package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Msg blabla
type Msg struct {
	cmd string
}

// Msg1 blabla
type Msg1 struct {
	Msg
	value string
}

// Msg2 blabla
type Msg2 struct {
	Msg
	value int
}

// GobDecode blabla
func (*Msg) GobDecode([]byte) error { return nil }

// GobEncode blabla
func (Msg) GobEncode() ([]byte, error) { return nil, nil }

func (m *Msg) encode() (bb bytes.Buffer) {
	enc := gob.NewEncoder(&bb)
	err := enc.Encode(m)
	if err != nil {
		log.Fatal("Cannot encode! err=", err)
	}
	return
}

func (m *Msg) decode(bb *bytes.Buffer) {
	dec := gob.NewDecoder(bb)
	err := dec.Decode(m)
	if err != nil {
		log.Fatal("Cannot decode Msg! err=", err)
	}
	return
}
