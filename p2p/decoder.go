package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(src io.Reader, dst interface{}) error
}

type GobDecoder struct{}

func NewGobDecoder() *GobDecoder {
	return &GobDecoder{}
}

func (g *GobDecoder) Decode(src io.Reader, dst interface{}) error {
	return gob.NewDecoder(src).Decode(dst)
}

type BuffDecoder struct{}

func NewBuffDecoder() *BuffDecoder {
	return &BuffDecoder{}
}

func (g *BuffDecoder) Decode(src io.Reader, data interface{}) error {
	buff := make([]byte, 1028)

	n, err := src.Read(buff)
	if err != nil {
		return err
	}
	data = buff[:n]
	return nil
}
