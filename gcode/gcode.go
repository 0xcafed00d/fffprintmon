package gcode

/*
 	Gcode command implemented:
	 M601 S1 - take control over wifi
	 M115    - Machine Info
	 M119    = get machine status
	 M27     - print status
	 M105    - read temperature
	 M28 [filesize] 0:/user/[filename]
	 		 - upload file
	 M29     - end file upload
	 M23 0:/user/[filename]
	 		 - print file

command flow:
	// get status
	M601 S1
	M115
wait_for_file:
	M119
	M105
	if no file goto wait_for_file
	// upload file
	m28
	m29
	// start print
	m23
print_loop:
	M119
	M105
	M27
	if Machine Status == BUILDING_FROM_SD goto print_loop

	// end



*/

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
