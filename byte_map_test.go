package main

import (
	"testing"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"math/rand"
	"os"
	"strconv"
	"github.com/cloudflare/ahocorasick"
)

var randomTargetSet []string
var randomTargetByteSet [][]byte
var randomTermSet []string
var termIndex *ByteMap
var dict *ahocorasick.Matcher

var TARGETS = 1000
var TARGET_LEN = 150
var TERMS = 100000
var TERM_LEN_MIN = 5
var TERM_LEN_MAX = 15

func init() {
	res, _ := strconv.Atoi(os.Getenv("TARGET_LEN"))
	if res != 0 {
		TARGET_LEN = res
	}
	res, _ = strconv.Atoi(os.Getenv("TERMS"))
	if res != 0 {
		TERMS = res
	}
	res, _ = strconv.Atoi(os.Getenv("TERM_LEN_MIN"))
	if res != 0 {
		TERM_LEN_MIN = res
	}
}

func assert(t *testing.T, b bool) {
	if !b {
		fmt.Printf("%+v\n", errors.New(""))
		t.Fatal()
	}
}

func TestGenerateIndex(t *testing.T) {
	index := GenerateIndex([]string{"foo", "baz", "bar🔥"})
	foo := "foo"
	assert(t, *index.node[foo[0]].node[foo[1]].node[foo[2]].result == "foo")
	baz := "baz"
	assert(t, *index.node[baz[0]].node[baz[1]].node[baz[2]].result == "baz")
	bar := "bar🔥"
	assert(t, *index.node[bar[0]].node[bar[1]].node[bar[2]].node[bar[3]].node[bar[4]].node[bar[5]].node[bar[6]].result == "bar🔥")
}

func TestIndexOf(t *testing.T) {
	index := GenerateIndex([]string{"foo", "baz", "bar🔥"})
	pos, match := index.Find("aaaafoo")
	assert(t, pos == 4 && *match == "foo")
	pos, match = index.Find("fobaaaaa")
	assert(t, pos == -1 && match == nil)
	pos, match = index.Find("マルチバイト文字")
	assert(t, pos == -1 && match == nil)
	pos, match = index.Find("マルチバイトbar🔥文字")
	assert(t, pos == 18 && *match == "bar🔥")
}

func InitBenchmark() {
	// 検索対象を初期化
	if randomTargetSet == nil {
		targetSet := make([]string, TARGETS)
		targetByteSet := make([][]byte, TARGETS)
		rand.Seed(0)
		for i := range targetSet {
			targetSet[i] = RandString(TARGET_LEN)
			targetByteSet[i] = []byte(targetSet[i])
		}
		randomTargetSet = targetSet
		randomTargetByteSet = targetByteSet
	}
	// ブラックリストを初期化
	if randomTermSet == nil {
		termSet := make([]string, TERMS)
		for i := range termSet {
			if TERM_LEN_MAX == TERM_LEN_MIN {
				termSet[i] = RandString(TERM_LEN_MIN)
			} else {
				termSet[i] = RandString(TERM_LEN_MIN + rand.Intn(TERM_LEN_MAX-TERM_LEN_MIN))
			}
		}
		randomTermSet = termSet
	}
	if termIndex == nil {
		index := GenerateIndex(randomTermSet)
		termIndex = &index
	}
	//if dict == nil {
	//	dict = ahocorasick.NewStringMatcher(randomTermSet)
	//}
}

func BenchmarkSearchIndexOf(b *testing.B) {
	InitBenchmark()
	inputArr := []string{"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36"}
	//dictArr := []string{
	//	"libcurl", "Wget", "Googlebot", "APIs-Google",
	//	"Bingbot", "Slurp", "DuckDuckBot", "Baiduspider",
	//	"YandexBot", "Sogou", "Exabot", "facebot",
	//	"facebookexternalhit", "ia_archiver", "Twitterbot", "Slackbot",
	//}
	dictArr := randomTermSet
	foundCount := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		input := inputArr[i%len(inputArr)]
		for _, term := range dictArr {
			res := strings.Index(input, term)
			if res != -1 {
				foundCount += 1
				break
			}
		}
	}
}

func BenchmarkSearchNormal(b *testing.B) {
	InitBenchmark()
	b.ResetTimer()
	foundCount := 0
	for i := 0; i < b.N; i++ {
		target := randomTargetSet[i%TARGETS]
		for t := range randomTermSet {
			term := randomTermSet[t]
			res := strings.Index(target, term)
			if res != -1 {
				foundCount += 1
				break
			}
		}
	}
}

func BenchmarkSearchBitMap(b *testing.B) {
	InitBenchmark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		termIndex.Find(randomTargetSet[i%TARGETS])
	}
}

func BenchmarkSearchAhocorasick(b *testing.B) {
	InitBenchmark()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Match([]byte(randomTargetByteSet[i%TARGETS]))
	}
}


type StringLike struct {
	data *[]byte
	len int
}