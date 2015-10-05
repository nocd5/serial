package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jacobsa/go-serial/serial"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	// required
	PortName string `short:"p" long:"port" description:"Serial Port" required:"true"`
	BaudRate uint   `short:"b" long:"baud" description:"Baud Rate" required:"true"`
	// optional
	DataBits   uint   `long:"data" description:"Number of Data Bits" default:"8"`
	ParityMode string `long:"parity" description:"Parity Mode. none/even/odd" default:"none"`
	StopBits   uint   `long:"stop" description:"Number of Stop Bits" default:"1"`
}

var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	var parityMode serial.ParityMode
	switch opts.ParityMode {
	case "none":
		parityMode = serial.PARITY_NONE
	case "odd":
		parityMode = serial.PARITY_ODD
	case "even":
		parityMode = serial.PARITY_EVEN
	default:
		fmt.Fprintf(os.Stderr, "Invalid ParityMode: %s\n", opts.ParityMode)
		fmt.Fprintf(os.Stderr, "`--parity` should be any one of none/odd/even\n")
		os.Exit(1)
	}

	options := serial.OpenOptions{
		PortName:   opts.PortName,
		BaudRate:   opts.BaudRate,
		DataBits:   opts.DataBits,
		ParityMode: parityMode,
		StopBits:   opts.StopBits,
	}

	port, err := serial.Open(options)
	defer port.Close()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := port.Write([]byte(scanner.Text()))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	buf := make([]byte, 128)
	done := make(chan bool)
	for {
		go func() {
			n, err := port.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			if n > 0 {
				fmt.Fprintf(os.Stdout, "%s", buf[:n])
			}
			done <- true
		}()
		<-done
	}
}
