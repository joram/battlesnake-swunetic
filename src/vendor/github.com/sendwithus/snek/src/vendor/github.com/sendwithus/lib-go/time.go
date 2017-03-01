package swu

import (
	"time"
)

func UnixMilliToNano(unixMilli int64) int64 {
	return unixMilli * 1000000
}

func UnixNanoToMilli(unixNano int64) int64 {
	return unixNano / 1000000
}

func TimeToUnixMilli(t time.Time) int64 {
	return UnixNanoToMilli(t.UnixNano())
}
