package benchmark

import (
	"time"
)

var (
	items = make(map[string]time.Time)
)

/**
计时开始
*/
func Start(name string) {
	items[name] = time.Now()
}

/**
结束计时
*/
func stopTime(name string) time.Duration {
	item, bl := items[name]
	if !bl {
		return time.Now().Sub(time.Now())
	}
	delete(items, name)

	return time.Now().Sub(item)
}

/**
结束时间, 并返回时间差(秒)
*/
func StopSeconds(name string) float64 {
	return stopTime(name).Seconds()
}

/**
结束时间, 并返回时间差(毫秒)
*/
func StopMilliseconds(name string) int64 {
	return stopTime(name).Milliseconds()
}

/**
结束时间, 并返回时间差(纳秒)
*/
func StopNanoseconds(name string) int64 {
	return stopTime(name).Nanoseconds()
}
