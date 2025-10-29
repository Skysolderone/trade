package utils

import (
	"fmt"
	"strconv"
)

func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return f
}
