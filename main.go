package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tarm/serial"
)

type Options struct {
	Port string `short:"p" long:"port" description:"Serial Port" required:"true"`
	Baud int    `short:"b" long:"baud" description:"Baud Rate" required:"true"`
}

var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	config := &serial.Config{Name: opts.Port, Baud: opts.Baud}
	port, err := serial.OpenPort(config)
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
