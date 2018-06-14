package main

import (
	"testing"
	"strconv"
	"os"
)

var NUM_ITEMS = 10

func init() {
	res, _ := strconv.Atoi(os.Getenv("NUM_ITEMS"))
	if res != 0 {
		NUM_ITEMS = res
	}

}

func BenchmarkMap(b *testing.B) {
	targetMap := make(map[int]int)
	for i := 0; i < NUM_ITEMS; i ++ {
		targetMap[i] = i
	}
	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		sum += targetMap[i % NUM_ITEMS]
	}
}

func BenchmarkArr(b *testing.B) {
	targetMap := make([]int, NUM_ITEMS)
	for i := 0; i < NUM_ITEMS; i ++ {
		targetMap[i] = i
	}
	b.ResetTimer()

	sum := 0
	for i := 0; i < b.N; i++ {
		sum += targetMap[i % NUM_ITEMS]
	}
}
