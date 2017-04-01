package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"time"

	"github.com/simulatedsimian/fffprintmon/gcode"
)

type Config struct {
	Help            bool
	PrinterHostName string
	PrinterPortNum  int
	MonitorFolder   string
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "display help")
	flag.IntVar(&config.PrinterPortNum, "p", 8899, "Printer Port number. Default: 8899")
	flag.StringVar(&config.PrinterHostName, "n", "", "Printer Host Name/IP address.")
	flag.StringVar(&config.PrinterHostName, "m", ".", "Folder to monitor for uploads. Default: Current folder")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: fffprintmon [options] ")
		flag.PrintDefaults()
	}
}

func exitOnError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v", message, err)
		os.Exit(1)
	}

}

func main() {
	flag.Parse()

	if config.PrinterHostName == "" || config.Help {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", config.PrinterHostName, config.PrinterPortNum))
	exitOnError(err, "Failed to connect to printer")
	defer conn.Close()

	gcode, err := gcode.Make(conn)
	exitOnError(err, "Error communicating with printer")

	run(gcode)
}

func run(gcode *gcode.GCode) {
	gcode.SendCommand("M601 S1")
	time.Sleep(5 * time.Second)
	gcode.SendCommand("M115")
	time.Sleep(5 * time.Second)
}
