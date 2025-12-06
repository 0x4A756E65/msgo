package msgo

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	maxDuration = time.Duration(1<<63 - 1)
	minDuration = -maxDuration - 1

	millisecond = time.Millisecond
	second      = time.Second
	minute      = time.Minute
	hour        = time.Hour
	day         = 24 * hour
	week        = 7 * day
	year        = day * 36525 / 100
	month       = year / 12
)

var (
	unitMap = map[string]time.Duration{
		"mo":           month,
		"month":        month,
		"months":       month,
		"ms":           millisecond,
		"msec":         millisecond,
		"msecs":        millisecond,
		"millisecond":  millisecond,
		"milliseconds": millisecond,
		"s":            second,
		"sec":          second,
		"secs":         second,
		"second":       second,
		"seconds":      second,
		"m":            minute,
		"min":          minute,
		"mins":         minute,
		"minute":       minute,
		"minutes":      minute,
		"h":            hour,
		"hr":           hour,
		"hrs":          hour,
		"hour":         hour,
		"hours":        hour,
		"d":            day,
		"day":          day,
		"days":         day,
		"w":            week,
		"week":         week,
		"weeks":        week,
		"y":            year,
		"yr":           year,
		"yrs":          year,
		"year":         year,
		"years":        year,
	}
	parseRe = regexp.MustCompile(`^\s*([+-]?(?:\d+)?\.?\d+)\s*([a-zA-Z]+)?\s*$`)
)

type unitDef struct {
	dur   time.Duration
	short string
	long  string
}

var orderedUnits = []unitDef{
	{dur: year, short: "y", long: "year"},
	{dur: month, short: "mo", long: "month"},
	{dur: week, short: "w", long: "week"},
	{dur: day, short: "d", long: "day"},
	{dur: hour, short: "h", long: "hour"},
	{dur: minute, short: "m", long: "minute"},
	{dur: second, short: "s", long: "second"},
	{dur: millisecond, short: "ms", long: "millisecond"},
}

// Parse converts a human string like "2h", "2 days", or "1y" into a time.Duration.
// If no unit is provided, the value is interpreted as milliseconds.
func Parse(s string) (time.Duration, error) {
	if len(s) == 0 || len(s) > 100 {
		return 0, fmt.Errorf("msgo: invalid length for %q", s)
	}

	match := parseRe.FindStringSubmatch(s)
	if match == nil {
		return 0, fmt.Errorf("msgo: invalid duration %q", s)
	}

	valueStr := match[1]
	unitStr := match[2]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("msgo: invalid number %q", valueStr)
	}

	unitKey := "ms"
	if unitStr != "" {
		unitKey = strings.ToLower(unitStr)
	}

	unit, ok := unitMap[unitKey]
	if !ok {
		return 0, fmt.Errorf("msgo: unknown unit %q", unitStr)
	}

	maxForUnit := float64(maxDuration) / float64(unit)
	minForUnit := float64(minDuration) / float64(unit)
	if value > maxForUnit || value < minForUnit {
		return 0, fmt.Errorf("msgo: duration overflow for %q", s)
	}

	duration := time.Duration(value * float64(unit))
	return duration, nil
}

// ParseStrict behaves the same as Parse; provided for API parity with Vercel/ms.
func ParseStrict(s string) (time.Duration, error) {
	return Parse(s)
}

// FormatShort renders a duration using the largest sensible unit, rounding to the nearest whole unit.
// Example: 90s => "2m".
func FormatShort(d time.Duration) string {
	return formatDuration(d, false)
}

// Format renders a duration in either the short or long style.
// This is a small convenience wrapper around FormatShort and FormatLong.
func Format(d time.Duration, long bool) string {
	if long {
		return FormatLong(d)
	}
	return FormatShort(d)
}

// FormatLong renders a duration using the largest sensible unit, rounding to the nearest whole unit.
// Example: 90s => "2 minutes".
func FormatLong(d time.Duration) string {
	return formatDuration(d, true)
}

func formatDuration(d time.Duration, long bool) string {
	if d == 0 {
		if long {
			return "0 ms"
		}
		return "0ms"
	}

	sign := ""
	abs := d
	if d < 0 {
		sign = "-"
		if d == minDuration {
			abs = maxDuration
		} else {
			abs = -d
		}
	}

	for i, u := range orderedUnits {
		last := i == len(orderedUnits)-1
		if abs < u.dur && !last {
			continue
		}

		count := math.Round(float64(abs) / float64(u.dur))
		if long {
			if u.short == "ms" {
				return fmt.Sprintf("%s%d ms", sign, int64(count))
			}
			name := u.long
			if count != 1 {
				name += "s"
			}
			return fmt.Sprintf("%s%d %s", sign, int64(count), name)
		}
		return fmt.Sprintf("%s%d%s", sign, int64(count), u.short)
	}

	if long {
		return "0 ms"
	}
	return "0ms"
}
