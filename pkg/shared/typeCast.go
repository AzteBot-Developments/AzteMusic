package shared

import "strconv"

func IdStringToInt(id string) int64 {
	num, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}
