package timezone

import "time"

func JakartaTz() *time.Location {
	tz, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.UTC
	}
	return tz
}
