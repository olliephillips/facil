package main

import (
	"log"
	"time"

	"github.com/olliephillips/facil/cmd"
)

func main() {
	t1 := time.Now()
	cmd.Execute()
	t2 := time.Now()
	log.Printf("Program ran in bout %v\n", t2.Sub(t1))
}
