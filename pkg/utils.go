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
