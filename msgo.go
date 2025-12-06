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

// Parse converts a human string like "2h", "2 days", or "1y" into a time.Duration.
// If no unit is provided, the value is interpreted as milliseconds.
func Parse(s string) (time.Duration, error) {
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

// FormatShort renders a duration using the largest sensible unit and short suffixes.
func FormatShort(d time.Duration) string {
	abs := absDuration(d)

	switch {
	case abs >= year:
		return fmt.Sprintf("%dy", roundJS(float64(d)/float64(year)))
	case abs >= month:
		return fmt.Sprintf("%dmo", roundJS(float64(d)/float64(month)))
	case abs >= week:
		return fmt.Sprintf("%dw", roundJS(float64(d)/float64(week)))
	case abs >= day:
		return fmt.Sprintf("%dd", roundJS(float64(d)/float64(day)))
	case abs >= hour:
		return fmt.Sprintf("%dh", roundJS(float64(d)/float64(hour)))
	case abs >= minute:
		return fmt.Sprintf("%dm", roundJS(float64(d)/float64(minute)))
	case abs >= second:
		return fmt.Sprintf("%ds", roundJS(float64(d)/float64(second)))
	default:
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
}

// FormatLong renders a duration using the largest sensible unit and long suffixes.
func FormatLong(d time.Duration) string {
	abs := absDuration(d)

	switch {
	case abs >= year:
		return formatLongUnit(d, abs, year, "year")
	case abs >= month:
		return formatLongUnit(d, abs, month, "month")
	case abs >= week:
		return formatLongUnit(d, abs, week, "week")
	case abs >= day:
		return formatLongUnit(d, abs, day, "day")
	case abs >= hour:
		return formatLongUnit(d, abs, hour, "hour")
	case abs >= minute:
		return formatLongUnit(d, abs, minute, "minute")
	case abs >= second:
		return formatLongUnit(d, abs, second, "second")
	default:
		return fmt.Sprintf("%d ms", d.Milliseconds())
	}
}

func formatLongUnit(d, abs, unit time.Duration, name string) string {
	rounded := roundJS(float64(d) / float64(unit))
	if abs >= unit+unit/2 {
		name += "s"
	}
	return fmt.Sprintf("%d %s", rounded, name)
}

func roundJS(v float64) int64 {
	return int64(math.Floor(v + 0.5))
}

func absDuration(d time.Duration) time.Duration {
	if d >= 0 {
		return d
	}
	if d == minDuration {
		return maxDuration
	}
	return -d
}
