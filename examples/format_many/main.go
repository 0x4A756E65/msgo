package main

import (
	"fmt"
	"time"

	"github.com/0x4A756E65/msgo"
)

func main() {
	values := []time.Duration{
		500 * time.Millisecond,
		1500 * time.Millisecond,
		90 * time.Second,
		45 * time.Minute,
		150 * time.Minute,
		2 * time.Hour,
		36 * time.Hour,
		7 * 24 * time.Hour,
		30 * 24 * time.Hour,
		365*24*time.Hour + 6*time.Hour, // ~1 year
		-1500 * time.Millisecond,
		-36 * time.Hour,
	}

	for _, d := range values {
		fmt.Printf("%-20v -> %-12s | %s\n",
			d,
			msgo.FormatShort(d),
			msgo.FormatLong(d),
		)
	}
}
