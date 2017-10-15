package gcode

import "fmt"
import "strings"

func extractCommand(s string) (cmd string) {
	fmt.Sscanf(s, "CMD %s", &cmd)
	return
}

type extracter struct {
	prefix      string
	extractFunc func(s string, c *CommandResponse)
}

var extracters = map[string][]extracter{
	//	CMD M119 Received.
	//	Endstop: X-max: 1 Y-max: 1 Z-max: 1
	//	Status: S:1 L:0 J:0 F:1
	//	MachineStatus: READY
	//	MoveMode: READY

	"M119": []extracter{
		{"MachineStatus", func(s string, c *CommandResponse) {
			var v string
			fmt.Sscanf(s, "MachineStatus: %s", &v)
			c.Params["MachineStatus"] = v
		}},
		{"MoveMode", func(s string, c *CommandResponse) {
			var v string
			fmt.Sscanf(s, "MoveMode: %s", &v)
			c.Params["MoveMode"] = v
		}},
		{"Endstop", func(s string, c *CommandResponse) {
			var xmax, ymax, zmax string
			fmt.Sscanf(s, "Endstop: X-max: %s Y-max: %s Z-max: %s", &xmax, &ymax, &zmax)
			c.Params["X-max"] = xmax
			c.Params["Y-max"] = ymax
			c.Params["Z-max"] = zmax
		}},
	},

	// CMD M115 Received.
	// Machine Type: Flashforge Finder
	// Machine Name: My Finder
	// Firmware: V1.5 20161014
	// SN: 628E8895
	// X: 140  Y: 140  Z: 140
	// Tool Count: 1

	"M115": []extracter{
		{"Machine Type", func(s string, c *CommandResponse) {
			p := strings.Split(s, ": ")
			c.Params[p[0]] = p[1]
		}},
		{"Firmware", func(s string, c *CommandResponse) {
			p := strings.Split(s, ": ")
			c.Params[p[0]] = p[1]
		}},
		{"X:", func(s string, c *CommandResponse) {
			var x, y, z string
			fmt.Sscanf(s, "X: %s  Y: %s  Z: %s", &x, &y, &z)
			c.Params["X"] = x
			c.Params["Y"] = y
			c.Params["Z"] = z
		}},
	},
}
