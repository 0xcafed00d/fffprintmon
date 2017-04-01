package gcode

import (
	"fmt"
	"io"
	"os"
)

type GCode struct {
	rw io.ReadWriter
}

func Make(rw io.ReadWriter) (*GCode, error) {
	g := &GCode{}
	g.rw = rw
	go g.responseReader()

	return g, nil
}

func (g *GCode) SendCommand(cmd string) error {
	_, err := fmt.Fprintf(g.rw, "~%s\r\n", cmd)
	return err
}

func (g *GCode) responseReader() {
	io.Copy(os.Stdout, g.rw)
}
