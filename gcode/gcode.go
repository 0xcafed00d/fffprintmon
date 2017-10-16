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
	"log"
	"strings"
)

type GCode struct {
	reader   *bufio.Reader
	writer   io.Writer
	respChan chan CommandResponse
}

type CommandResponse struct {
	Params map[string]string
	err    error
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
	log.Printf("~%s", cmd)
	if err != nil {
		return CommandResponse{}, err
	}
	resp := <-g.respChan
	return resp, resp.err
}

func (g *GCode) CMDTakeControl() (CommandResponse, error) {
	return g.SendCommand("M601 S1")
}

func (g *GCode) CMDHomePos() (CommandResponse, error) {
	return g.SendCommand("G28")
}

func (g *GCode) CMDPrinterInfo() (CommandResponse, error) {
	return g.SendCommand("M115")
}

func (g *GCode) CMDPrinterStatus() (CommandResponse, error) {
	return g.SendCommand("M119")
}

func (g *GCode) CMDSetRGBLights(r, gr, b int) (CommandResponse, error) {
	return g.SendCommand(fmt.Sprintf("M146 r%d g%d b%d", r, gr, b))
}

func (g *GCode) CMDCoordAbs() (CommandResponse, error) {
	return g.SendCommand("G90")
}

func (g *GCode) CMDCoordRel() (CommandResponse, error) {
	return g.SendCommand("G91")
}

func (g *GCode) CMDMoveXYZ(x, y, z float64) (CommandResponse, error) {
	return g.SendCommand(fmt.Sprintf("G1 X%f Y%f Z%f", x, y, z))
}

func (g *GCode) CMDGetXYZ() (CommandResponse, error) {
	return g.SendCommand("M114")
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
				resp.Params["CMD"] = extractCommand(line)
			} else {
				e := extracters[resp.Params["CMD"]]
				if len(e) > 0 {
					for i := range e {
						if strings.HasPrefix(line, e[i].prefix) {
							e[i].extractFunc(line, &resp)
						}
					}
				}
			}
		}
	}
}
