package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	keywordFile = "./keywords"
	matchTexts  = `分割侵权行为电站证券市场营业执照违约全球化市盈率推迟虐待增值税违法行为未成年人公司董事会没收有车有效期`
)

func getMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func TestRegExp(t *testing.T) {
	fp, err := os.Open(keywordFile)
	if err != nil {
		return
	}

	defer func() { _ = fp.Close() }()

	var regexpFmt string

	buf := bufio.NewReader(fp)
	for {
		word, err := buf.ReadString('\n')
		if err != nil {
			break
		}

		word = strings.TrimSpace(word)
		if len(word) > 0 {
			regexpFmt += word + "|"
		}
	}

	regexpFmt = strings.Trim(regexpFmt, "|")

	matcher := regexp.MustCompile(regexpFmt)

	startTime := getMillisecond()

	matchs := matcher.FindAllString(matchTexts, -1)

	fmt.Printf("millisecond: %d\n", getMillisecond()-startTime)
	fmt.Println(len(matchs))

	for _, match := range matchs {
		fmt.Println(match)
	}
}

func TestDict(t *testing.T) {
	dict := NewDict()

	dict.LoadWordsFile(keywordFile)

	texts := []byte(matchTexts)

	startTime := getMillisecond()

	matchs := dict.MatchAll(texts, -1)

	fmt.Printf("millisecond: %d\n", getMillisecond()-startTime)
	fmt.Println(len(matchs))

	for _, match := range matchs {
		fmt.Println(match)
	}
}
