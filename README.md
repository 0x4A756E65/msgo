# msgo

![msgo banner](assets/ms-go-banner.png)

Go clone of Vercel’s `ms`: parse human-friendly duration strings to `time.Duration` and format durations back into short or long strings.

## Install

```sh
go get github.com/0x4A756E65/msgo
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/0x4A756E65/msgo"
)

func main() {
	// Parse
	d, err := msgo.Parse("2.5 hrs")
	if err != nil {
		panic(err)
	}
	fmt.Println(d) // 2h30m0s

	// Format
	fmt.Println(msgo.FormatShort(90 * time.Second)) // "2m"
	fmt.Println(msgo.FormatLong(2 * time.Hour))     // "2 hours"
	fmt.Println(msgo.Format(36*time.Hour, true))    // "2 days"
}
```

## API

- `Parse` / `ParseStrict`: parse strings like `"2 days"`, `"1.5h"`, or `"-.5w"` into `time.Duration` (milliseconds default).
- `FormatShort`: short output that mirrors `ms(value)`/`format(value)` (`90s` → `2m`, `1500ms` → `2s`).
- `FormatLong`: long output that mirrors `{ long: true }` (`1500ms` → `2 seconds`).
- `Format(d, long)`: convenience wrapper to flip between short and long output in one call.

## Supported input units

Milliseconds (`ms`, `msec`, `msecs`, `millisecond`, `milliseconds`), seconds, minutes, hours, days, weeks, months (`mo`), and years (`y`). Input is case-insensitive, can include spaces (`"2 days"` or `"2days"`), and supports decimals/negatives (`"1.5h"`, `"-0.5d"`, `".5ms"`). Bare numbers default to milliseconds.

## Formatting rules

- Short/Long choose the largest sensible unit and round to the nearest whole unit (`90s` → `2m`; `5400s` → `2 hours`; `1500ms` → `2 seconds`).
- Values smaller than a millisecond round to `0ms`/`0 ms` because `time.Duration` is integer-based.

## License

MIT
