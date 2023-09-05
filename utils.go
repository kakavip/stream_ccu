package main

import "math/rand"

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func uniqueArray(array []string) []string {
	tempMap := make(map[string]bool)
	for _, v := range array {
		tempMap[v] = true
	}

	result := make([]string, 0)
	for k, _ := range tempMap {
		result = append(result, k)
	}
	return result
}

func RemoveItem[T string | int](arr []T, item T) []T {
	var new []T
	index := 0
	for _, i := range arr {
		if i != item {
			new[index] = i
			index++
		}
	}
	return new[:index]
}
