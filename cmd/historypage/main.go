package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/thecubic/medtronic"
)

func main() {
	pump := medtronic.Open()
	defer pump.Close()
	switch len(os.Args) {
	case 1:
		result := pump.HistoryPageCount()
		if pump.Error() != nil {
			log.Fatal(pump.Error())
		}
		fmt.Println(result)
	case 2:
		page, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		data := pump.HistoryPage(page)
		if pump.Error() != nil {
			log.Fatal(pump.Error())
		}
		fmt.Printf("% X\n", data)
	default:
		log.Fatalf("Usage: %s [page#]", os.Args[0])
	}
}
