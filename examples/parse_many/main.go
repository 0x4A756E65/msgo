package main

import (
	"fmt"
	"log"
	"time"

	"github.com/0x4A756E65/msgo"
)

func main() {
	samples := []string{
		"100",
		"250ms",
		"1s",
		"1.5s",
		"45m",
		"2h",
		"1d",
		"1w",
		"1mo",
		"1y",
		"-500ms",
		"-.5h",
		"3 DAYS",
		"0.5 minutes",
	}

	for _, raw := range samples {
		d, err := msgo.Parse(raw)
		if err != nil {
			log.Printf("%-12q -> error: %v", raw, err)
			continue
		}
		fmt.Printf("%-12q -> %-12s | %-18s | %v\n",
			raw,
			msgo.FormatShort(d),
			msgo.FormatLong(d),
			d,
		)
	}

	// Show how durations can be used directly.
	d, _ := msgo.Parse("1.5s")
	fmt.Println()
	fmt.Println("Sleeping for", d)
	time.Sleep(d)
	fmt.Println("Done.")
}
