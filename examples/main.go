package main

import (
	"fmt"
	"log"
	"time"

	"github.com/0x4A756E65/msgo"
)

func main() {
	raw := "2.5s"

	interval, err := msgo.Parse(raw)
	if err != nil {
		log.Fatalf("failed to parse %q: %v", raw, err)
	}

	fmt.Printf("parsed interval:  %v\n", interval)
	fmt.Printf("as milliseconds: %dms\n", interval.Milliseconds())
	fmt.Printf("short:           %s\n", msgo.FormatShort(interval))
	fmt.Printf("long:            %s\n", msgo.FormatLong(interval))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("Ticking 3 times...")
	for i := range 3 {
		<-ticker.C
		fmt.Printf("tick %d at %s\n", i+1, time.Now().Format(time.RFC3339))
	}
}
