package util

import "time"

type Clock struct{}

func NewClock() *Clock {
	return &Clock{}
}

// GetRealTimeInNSecs returns the nano seconds from epoch.
func (c *Clock) GetRealTimeInNSecs() uint64 {
	return uint64(time.Now().Unix())
}

// GetOSTime returns the OS time.
func (c *Clock) GetOSTime(ty int) int {
	now := time.Now()
	switch ty {
	case 1:
		return now.Second()
	case 2:
		return now.Minute()
	case 3:
		return now.Hour()
	case 4:
		return now.Day()
	case 5:
		return int(now.Month())
	case 6:
		return now.Year()
	case 7:
		return int(now.Weekday())
	default:
		return now.Second()
	}
}

// NanoToMicro converts nanoseconds to microseconds.
func NanoToMicro(x uint64) uint64 {
	return x / 1000
}

// NanoToMill converts nanoseconds to millseconds.
func NanoToMill(x uint64) uint64 {
	return NanoToMicro(x) / 1000
}
