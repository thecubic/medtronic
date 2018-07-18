package main

import (
	"log"
	"time"

	"github.com/thecubic/medtronic"
	"github.com/thecubic/medtronic/packet"
)

const (
	verbose = true
)

func main() {
	if verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	}
	pump := medtronic.Open()
	defer pump.Close()
	for pump.Error() == nil {
		p, rssi := pump.Radio.Receive(time.Hour)
		if pump.Error() != nil {
			log.Print(pump.Error())
			pump.SetError(nil)
			continue
		}
		if verbose {
			log.Printf("raw data: % X (RSSI = %d)", p, rssi)
		}
		data, err := packet.Decode(p)
		if err != nil {
			log.Print(err)
			continue
		}
		if verbose {
			log.Printf("decoded:  % X", data)
		} else {
			log.Printf("% X (RSSI = %d)", data, rssi)
		}

	}
	log.Fatal(pump.Error())
}
