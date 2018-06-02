package main

import (
	"math/rand"
	"fmt"
	"time"
	"strings"
)

func init() {
}

var strArray = []rune("abcdefghijklmnopqrstuvwxyz")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = strArray[rand.Intn(len(strArray))]
	}
	return string(b)
}

type ByteMap struct {
	values  [256]*ByteMap
	results [256]*string
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
	blackIndex := ByteMap{[256]*ByteMap{}, [256]*string{}}
	for bi := range blackSet {
		blackTerm := blackSet[bi]
		currentIndex := &blackIndex
		for i := 0; i < len(blackTerm)-1; i++ {
			code := blackTerm[i]
			if currentIndex.values[code] == nil {
				currentIndex.values[code] = &ByteMap{[256]*ByteMap{}, [256]*string{}}
			}
			currentIndex = currentIndex.values[code]
		}
		lastCode := blackTerm[len(blackTerm)-1]
		currentIndex.results[lastCode] = &blackTerm
	}
	//fmt.Printf("  %v\n", *blackIndex.values[117].values[106].values[113].values[106].results[100]) // ujqjd
	fmt.Printf("ImplimentB index done time: %f\n", time.Since(start).Seconds())

	foundCount = 0
	for ti := range targetSet {
		target := targetSet[ti]
		found := -1 // 検索対象が存在した位置
	SEARCH_LOOP:
		for i := 0; i < len(target); i++ { // targetの文字ごとにループを回す
			currentIndex := &blackIndex
			for i2 := 0; i+i2 < len(target); i2++ {
				code := target[i+i2]
				if currentIndex.results[code] != nil {
					found = i + i2
					break SEARCH_LOOP // 検索対象が含まれることが分かったので、検索終了
				}
				if currentIndex.values[code] == nil { // Indexに存在しないので検索終了
					break
				}
				currentIndex = currentIndex.values[code] // 次のTreeは見つかったが、resultは無いのでループ続行
			}
		}
		if found != -1 {
			foundCount += 1
		}
	}
	fmt.Printf("ImplimentB founds: %d time: %f\n", foundCount, time.Since(start).Seconds())
}
