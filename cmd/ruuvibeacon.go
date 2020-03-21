package main

import (
	"github.com/peknur/ruuvibeacon"
	_ "github.com/peknur/ruuvibeacon/publishers"
)

func main() {
	ruuvibeacon.Run()
}
