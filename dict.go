package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/liexusong/cedar-go"
)

type Dict struct {
	lastId int
	dTree  *cedar.Cedar
	cache  map[int]string
	mutex  *sync.RWMutex
}

func NewDict() *Dict {
	return &Dict{
		lastId: 0,
		dTree:  cedar.New(),
		cache:  make(map[int]string),
		mutex:  &sync.RWMutex{},
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

func (dict *Dict) AddKeyword(realText string) bool {
	byteText := []byte(realText)

	dict.mutex.Lock()
	defer dict.mutex.Unlock()

	_, err := dict.dTree.Get(byteText)
	if err == nil {
		return false
	}

	lastId := dict.GenLastId()

	if err = dict.dTree.Insert(byteText, lastId); err != nil {
		return false
	}

	dict.cache[lastId] = realText

	return true
}

func (dict *Dict) DelKeyword(realText string) bool {
	byteText := []byte(realText)

	dict.mutex.Lock()
	defer dict.mutex.Unlock()

	id, err := dict.dTree.Get(byteText)
	if err != nil {
		return false
	}

	if err = dict.dTree.Delete(byteText); err != nil {
		return false
	}

	delete(dict.cache, id)

	return true
}

func (dict *Dict) MatchAll(text []byte, size int) []string {
	var (
		values []string
	)

	dict.mutex.RLock()
	defer dict.mutex.RUnlock()

	matches := dict.dTree.MatchAll(text, size)
	if len(matches) > 0 {
		for _, match := range matches {
			if value, exists := dict.cache[match]; exists {
				values = append(values, value)
			}
		}
	}

	return values
}

func (dict *Dict) Exists(text []byte) bool {
	dict.mutex.RLock()
	defer dict.mutex.RUnlock()

	return dict.dTree.Exists(text)
}

func (dict *Dict) LoadWordsFromFile(path string) error {
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
			dict.AddKeyword(word)
		}
	}
}
