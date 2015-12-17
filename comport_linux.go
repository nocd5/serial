package main

import (
	"fmt"
	"os"
)

func listComPorts() {
	fmt.Fprintln(os.Stderr, "the flag -l, --list are available on windows only.")
}
