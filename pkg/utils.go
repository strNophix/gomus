package gomus

import "log"

func mapList[T any, R any](l []T, f func(T) R) []R {
	var c []R = []R{}
	for _, item := range l {
		c = append(c, f(item))
	}
	return c
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Numeric interface {
	float64
}

func MapFloatBetween[T Numeric](value, start1, stop1, start2, stop2 T) T {
	return start2 + (value-start1)*(stop2-start2)/(stop1-start1)
}
