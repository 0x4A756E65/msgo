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
fmt.Println(msgo.FormatShort(90 * time.Second)) // "1m 30s"
fmt.Println(msgo.FormatLong(2 * time.Hour))     // "2 hours"
fmt.Println(msgo.Format(36*time.Hour, true))    // "1 day 12 hours"
}
```

## Supported input units

Milliseconds (`ms`, `msec`, `msecs`, `millisecond`, `milliseconds`), seconds, minutes, hours, days, weeks, months (`mo`), and years (`y`). Input is case-insensitive, can include spaces (`"2 days"` or `"2days"`), and supports decimals/negatives (`"1.5h"`, `"-0.5d"`, `".5ms"`). Bare numbers default to milliseconds.

## Formatting rules

- Short/Long choose the largest sensible unit and include a remainder in the next smallest unit (`90s` → `1m 30s`; `2h30m` → `2 hours 30 minutes`).
- Values smaller than a millisecond print as `0ms`/`0 milliseconds` because `time.Duration` is integer-based.

## Notes vs Vercel/ms

- Go API: separate `Parse`, `FormatShort`, `FormatLong`, and a convenience `Format(d, long)`; no type-based overloading or `parseStrict`.
- Errors are returned instead of throwing/returning `NaN`.
- Input length guard: 1–100 characters.

## License

MIT
