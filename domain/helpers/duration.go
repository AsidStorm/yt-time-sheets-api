package helpers

import (
	"fmt"
	"github.com/sosodev/duration"
	"strings"
	"time"
)

var patterns = map[string]bool{
	"W": true,
	"D": true,
	"H": false,
	"M": false,
	"S": false,
}

func DurationFromString(in string) (string, error) {
	// 1w 2d 3h 4m 5s
	// P1W2DT3H4M5S
	parts := strings.Split(in, " ")

	var dateParts []string
	var timeParts []string

	counter := 0

	exists := make(map[string]bool)

	for pattern, place := range patterns {
		for _, part := range parts {
			value := strings.ToUpper(strings.TrimSpace(part))

			if strings.HasSuffix(value, pattern) {
				if exists[pattern] {
					return "", fmt.Errorf("doubled pattern [%s]", pattern)
				}

				exists[pattern] = true

				if place {
					dateParts = append(dateParts, value)
				} else {
					timeParts = append(timeParts, value)
				}

				counter++

				break
			}
		}
	}

	if counter != len(parts) {
		return "", fmt.Errorf("unexpected count [%d] of parts [%d]", counter, len(parts))
	}

	if len(dateParts) == 0 && len(timeParts) == 0 {
		return "", fmt.Errorf("unable to parse duration [%s] due [no parts]", in)
	}

	out := "P"

	if len(dateParts) > 0 {
		out += strings.Join(dateParts, "")
	}

	if len(timeParts) > 0 {
		out += "T"
		out += strings.Join(timeParts, "")
	}

	return out, nil
}

func ParseDuration(in string) (time.Duration, error) {
	lowerThanZero := strings.Contains(in, "-")

	d, err := duration.Parse(strings.ReplaceAll(in, "-", ""))
	if err != nil {
		return time.Second, err
	}

	// В Yandex Tracker'e 1W = 5D, а 1D = 8H, таким образом мы заменяем

	if d.Weeks > 0 {
		d.Days += d.Weeks * 5
		d.Weeks = 0
	}

	if d.Days > 0 {
		d.Hours += d.Days * 8
		d.Days = 0
	}

	if lowerThanZero {
		return 0 - d.ToTimeDuration(), nil
	}

	return d.ToTimeDuration(), nil
}
