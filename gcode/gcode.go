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
	"bufio"
	"fmt"
	"io"
	"strings"
)

type GCode struct {
	reader   *bufio.Reader
	writer   io.Writer
	respChan chan CommandResponse
}

type CommandResponse struct {
	Command string
	Params  map[string]string
	err     error
}

func MakeCommandResponse() CommandResponse {
	return CommandResponse{Params: make(map[string]string)}
}

func New(rw io.ReadWriter) *GCode {
	g := &GCode{}
	g.reader = bufio.NewReader(rw)
	g.writer = rw
	g.respChan = make(chan CommandResponse)
	go g.responseReader()
	return g
}

func (g *GCode) SendCommand(cmd string) (CommandResponse, error) {
	_, err := fmt.Fprintf(g.writer, "~%s\r\n", cmd)
	if err != nil {
		return CommandResponse{}, err
	}
	resp := <-g.respChan
	return resp, resp.err
}

func (g *GCode) responseReader() {
	resp := MakeCommandResponse()

	for {
		line, err := g.reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			g.respChan <- CommandResponse{err: err}
			return
		}

		if line == "ok" {
			g.respChan <- resp
			resp = MakeCommandResponse()
		} else {
			if strings.HasPrefix(line, "CMD") {
				resp.Command = line
			} else {
				i := strings.Index(line, ": ")
				if i != -1 {
					resp.Params[line[:i]] = line[i:]
				}
			}
		}
	}
}
