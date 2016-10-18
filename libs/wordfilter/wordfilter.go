package Wordfilter

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/huichen/sego"
)

var segmenter sego.Segmenter
var wordfilterMap = make(map[string]bool)

var gopath string = os.Getenv("GOPATH")

func init() {
	segmenter.LoadDictionary(gopath + "/src/github.com/k4s/tea/libs/wordfilter/dictionary.txt")
	wordlist := loadword()
	for _, v := range wordlist {
		wordfilterMap[v] = true
	}
}

//对外过滤接口
func Wordfilter(words string) string {
	text := []byte(words)
	segments := segmenter.Segment(text)
	wordslice := sego.SegmentsToSlice(segments, true)
	for _, word := range wordslice {
		if _, ok := wordfilterMap[word]; ok {
			words = strings.Replace(words, word, stringToStar(word), -1)
		}
	}
	return words
}

//加载敏感词
func loadword() []string {
	f, err := os.Open(gopath + "/src/github.com/k4s/tea/libs/wordfilter/wordfilter.txt")
	if err != nil {
		log.Fatal("%v", err)
	}
	defer f.Close()
	fd, err := ioutil.ReadAll(f)
	wordlist := strings.Split(string(fd), ",")
	return wordlist
}

//文字转换为"*"
func stringToStar(word string) string {
	var star string
	for i := 0; i < len(word)/3; i++ {
		star += "*"
	}
	return star
}
