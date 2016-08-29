package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	RESULTS = `tmp/1, 70c881d4a26984ddce795f6f71817c9cf4480e79, 4
tmp/2, 8aed1322e5450badb078e1fb60a817a1df25a2ca, 4
tmp/4/1, 01464e1616e3fdd5c60c0cc5516c1d1454cc4185, 4
`
)

func createDir(name string) {
	os.MkdirAll(name, 0777)
}

func removeDir(name string) {
	os.RemoveAll(name)
}

func createFile(name string, content string) {
	createDir(filepath.Dir(name))
	ioutil.WriteFile(name, []byte(content), 0644)
}

func TestSHA(t *testing.T) {
	defer removeDir("./tmp")
	createFile("./tmp/1", "aaaa")
	createFile("./tmp/2", "bbbb")
	createFile("./tmp/3", "cccc")
	createFile("./tmp/4/1", "dddd")
	createFile("./tmp/5/1", "eeee")
	results := walk("./tmp", []string{"tmp/3", "tmp/5/*"})
	if strings.Join(results, "") != RESULTS {
		t.Fail()
	}
}

// TODO
// - Performance Testing
// - Symbol Link Testing
