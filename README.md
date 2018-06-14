大量の単語が含まれる辞書から、最初にマッチする文字を検索する。

```go
type ByteMap struct {
	values  [256]*ByteMap
	results [256]*string
}
```

- "abcd" のASCIIコードは 97, 98, 99, 100
- `index.values[97].values[98].values[99].values[100].result = &"abcd"` となるようにIndexを構築する
- 1文字づつシフトして検索し、 "abcd" "bcd" "cd" "d" と順に検索する

```
> go run .\main.go
TARGETS: 1000 TARGET_LEN: 150 BLACKLIST: 100000 BLACK_TERM_LEN: 5
ImplimentA founds: 709 time: 6.439321
ImplimentB index done time: 0.329813
ImplimentB founds: 709 time: 0.341806
```

ベンチマーク

```
TERMS=10000 go test -bench . -benchtime 5s -benchmem
BenchmarkNormal-8                  10000           1146952 ns/op               0 B/op          0 allocs/op
BenchmarkBitMap-8                2000000              3602 ns/op               0 B/op          0 allocs/op
BenchmarkAhocorasick-8           1000000              9222 ns/op             160 B/op          1 allocs/op
PASS
ok      str-match-example-2018   35.663s
```
