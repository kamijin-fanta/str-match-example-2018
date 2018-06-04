package main

type ByteMap struct {
	node   [256]*ByteMap
	result *string
}

func (byteMap *ByteMap) Find(target string) (int, *string) {
	for i := 0; i < len(target); i++ { // targetの文字ごとにループを回す
		currentIndex := byteMap
		for i2 := 0; i+i2 < len(target); i2++ {
			code := target[i+i2]
			if currentIndex.node[code] == nil { // Indexに存在しないので検索終了
				break
			}
			if currentIndex.node[code].result != nil {
				return i, currentIndex.node[code].result // 検索対象が含まれることが分かったので、検索終了
			}
			currentIndex = currentIndex.node[code] // 次のTreeは見つかったが、resultは無いのでループ続行
		}
	}
	return -1, nil
}

func GenerateIndex(terms []string) ByteMap {
	index := ByteMap{[256]*ByteMap{}, nil}
	for bi := range terms {
		blackTerm := terms[bi]
		currentIndex := &index
		for i := 0; i < len(blackTerm); i++ {
			code := blackTerm[i]
			if currentIndex.node[code] == nil {
				currentIndex.node[code] = &ByteMap{[256]*ByteMap{}, nil}
			}
			currentIndex = currentIndex.node[code]
		}
		currentIndex.result = &blackTerm
	}
	//fmt.Printf("  %v\n", *index.node[117].node[106].node[113].node[106].results[100]) // ujqjd
	return index
}
