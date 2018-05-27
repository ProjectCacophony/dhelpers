package humanize

import (
	"strconv"
	"time"
)

// Duration formats a time.Duration in a human readable format
func Duration(d time.Duration) (result string) {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) - (hours * 60)
	seconds := int(d.Seconds()) - (minutes * 60) - (hours * 60 * 60)

	if hours > 0 {
		days := hours / 24
		hoursLeft := hours % 24
		if days > 0 {
			result += strconv.Itoa(days) + "d"
		}
		if hoursLeft > 0 {
			result += strconv.Itoa(hoursLeft) + "h"
		}
	}
	if minutes > 0 {
		result += strconv.Itoa(minutes) + "m"
	}
	if seconds > 0 {
		result += strconv.Itoa(seconds) + "s"
	}
	return result
}
