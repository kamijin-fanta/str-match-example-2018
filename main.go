package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var strArray = []rune("abcdefghijklmnopqrstuvwxyz")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = strArray[rand.Intn(len(strArray))]
	}
	return string(b)
}

func main() {
	TARGETS := 1000
	TARGET_LEN := 150
	BLACKLIST := 100000
	BLACK_TERM_LEN := 5

	// 検索対象を初期化
	targetSet := make([]string, TARGETS)
	rand.Seed(0)
	for i := range targetSet {
		targetSet[i] = RandString(TARGET_LEN)
	}

	// ブラックリストを初期化
	blackSet := make([]string, BLACKLIST)
	for i := range blackSet {
		blackSet[i] = RandString(BLACK_TERM_LEN)
	}

	fmt.Printf("TARGETS: %d TARGET_LEN: %d BLACKLIST: %d BLACK_TERM_LEN: %d\n", TARGETS, TARGET_LEN, BLACKLIST, BLACK_TERM_LEN)

	// 通常バージョン
	start := time.Now()
	foundCount := 0
	for ti := range targetSet {
		target := targetSet[ti]
		for bi := range blackSet {
			blackTerm := blackSet[bi]
			res := strings.Index(target, blackTerm)
			if res != -1 {
				//fmt.Printf("%s %d\n", blackTerm, res)
				foundCount += 1
				break
			}
		}
	}

	fmt.Printf("ImplimentA founds: %d time: %f\n", foundCount, time.Since(start).Seconds())

	// ByteMapバージョン
	start = time.Now()
	// Index構築
	blackIndex := GenerateIndex(blackSet)
	fmt.Printf("ImplimentB index done time: %f\n", time.Since(start).Seconds())

	foundCount = 0
	for ti := range targetSet {
		target := targetSet[ti]
		res, _ := blackIndex.Find(target)
		if res != -1 {
			foundCount += 1
		}
	}
	fmt.Printf("ImplimentB founds: %d time: %f\n", foundCount, time.Since(start).Seconds())
}
