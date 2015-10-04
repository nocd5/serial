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
	Port string `short:"p" long:"port" description:"Serial Port"`
	Baud int    `short:"b" long:"baud" description:"Baud Rate"`
}

var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	c := &serial.Config{Name: opts.Port, Baud: opts.Baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	buf := make([]byte, 256)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		go func() {
			for scanner.Scan() {
				_, err := s.Write([]byte(scanner.Text()))
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

		go func() {
			n, err := s.Read(buf)
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
