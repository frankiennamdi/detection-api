package services

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

func AddTime(initialTime int64, unit int, duration time.Duration) int64 {
	initial := time.Unix(initialTime, 0)
	return initial.Add(time.Duration(unit) * duration).Unix()
}

func AssertThatJSONEqual(s1, s2 string) (bool, error) {
	var o1 interface{}

	var o2 interface{}

	var err error

	err = json.Unmarshal([]byte(s1), &o1)

	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}

	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
