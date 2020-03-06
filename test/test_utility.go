package test

import (
	"time"
)

func AddTime(initialTime int64, unit int, duration time.Duration) int64 {

	initial := time.Unix(initialTime, 0)
	return initial.Add(time.Duration(unit) * duration).Unix()
}
