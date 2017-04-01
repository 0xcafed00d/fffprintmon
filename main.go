package main

import (
	"flag"
	"fmt"
	"net"
	"os"
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
	flag.StringVar(&config.PrinterHostName, "h", "", "Printer Host Name/IP address.")
	flag.StringVar(&config.PrinterHostName, "m", ".", "Folder to monitor for uploads. Default: Current folder")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: fffprintmon [options] ")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 && config.Help {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.Dial()

}
