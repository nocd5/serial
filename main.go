package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	// required
	PortName string `short:"p" long:"port" description:"Serial Port"`
	BaudRate uint   `short:"b" long:"baud" description:"Baud Rate"`
	// optional
	DataBits     uint   `long:"data" description:"Number of Data Bits" default:"8"`
	ParityMode   string `long:"parity" description:"Parity Mode. none/even/odd" default:"none"`
	StopBits     uint   `long:"stop" description:"Number of Stop Bits" default:"1"`
	ListComPorts bool   `short:"l" long:"list" description:"List COM Ports"`
	BinaryMode   bool   `short:"y" long:"binary" description:"Binary Mode"`
}

var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.ListComPorts {
		listComPorts()
		os.Exit(0)
	}

	if opts.PortName == "" || opts.BaudRate == 0 {
		fmt.Fprintln(os.Stderr, "the required flags `/b, /baud' and `/p, /port' were not specified")
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
			if opts.BinaryMode {
				if _, err := port.Write(getBytes(scanner.Text())); err != nil {
					log.Fatal(err)
				}
			} else {
				if _, err := port.Write([]byte(scanner.Text())); err != nil {
					log.Fatal(err)
				}
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

func getBytes(src string) []byte {
	var result []byte
	bytes := strings.Split(src, " ")
	for i := range bytes {
		val, _ := strconv.ParseInt(bytes[i], 0, 0)
		for 0 < val {
			result = append(result, byte(val&0xFF))
			val = val >> 8
		}
	}
	return result
}
