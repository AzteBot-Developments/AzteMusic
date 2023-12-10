package shared

import (
	"fmt"
	"time"
)

// Formats the passed in duration (given as seconds) into a nice format and returns the result.
func FormatDuration(seconds int64) string {

	duration := time.Second * time.Duration(seconds)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	secondsS := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secondsS)
	}

	return fmt.Sprintf("%02d:%02d", minutes, secondsS)
}
