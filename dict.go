package main

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/liexusong/cedar-go"
)

type Dict struct {
	lastId   int
	wordTree *cedar.Cedar
	wordMaps map[int]string
}

func NewDict() *Dict {
	return &Dict{
		lastId:   0,
		wordTree: cedar.New(),
		wordMaps: make(map[int]string),
	}
}

func (dict *Dict) GenLastId() int {
	lastId := dict.lastId
	dict.lastId++
	return lastId
}

func (dict *Dict) GetLastId() int {
	return dict.lastId
}

func (dict *Dict) AddWord(origText string) bool {
	realText := []byte(origText)

	_, err := dict.wordTree.Get(realText) // Find the word aready exists?
	if err == nil {
		return false
	}

	lastId := dict.GenLastId()

	if err = dict.wordTree.Insert(realText, lastId); err != nil {
		return false
	}

	dict.wordMaps[lastId] = origText

	return true
}

func (dict *Dict) DelWord(origText string) bool {
	realText := []byte(origText)

	index, err := dict.wordTree.Get(realText) // Find the word aready exists?
	if err != nil {
		return false
	}

	if err = dict.wordTree.Delete(realText); err != nil {
		return false
	}

	delete(dict.wordMaps, index)

	return true
}

func (dict *Dict) MatchAll(text []byte, size int) []string {
	var values []string

	matchs := dict.wordTree.MatchAll(text, size)
	if len(matchs) > 0 {
		for _, match := range matchs {
			if value, exists := dict.wordMaps[match]; exists {
				values = append(values, value)
			}
		}
	}

	return values
}

func (dict *Dict) LoadWordsFile(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() { _ = fp.Close() }()

	buf := bufio.NewReader(fp)
	for {
		word, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		word = strings.TrimSpace(word)
		if len(word) > 0 {
			dict.AddWord(word)
		}
	}

	return nil
}
