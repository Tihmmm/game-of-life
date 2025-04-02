package util

import "math/rand"

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CreateRandomUintSlice(size int, max uint) []uint {
	var result []uint

	for len(result) < size {
		num := rand.Intn(int(max))
		result = append(result, uint(num))
	}

	return result
}
