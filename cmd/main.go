package main

import (
	"flag"
	"github.com/jkittell/mediastreaminfo"
	"log"
)

func main() {
	port := flag.Int("port", 3000, "port number")
	debug := flag.Bool("debug", false, "debugging")
	flag.Parse()

	mediastreaminfo.Debug(*debug)
	err := mediastreaminfo.StartService(*port)
	if err != nil {
		log.Println(err)
	}
}
