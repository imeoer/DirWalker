package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"github.com/imeoer/dirwalker/limiter"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ignore(patterns []string, path string) bool {
	for _, pattern := range patterns {
		if matched, err := filepath.Match(pattern, path); err == nil && matched {
			return true
		}
	}
	return false
}

func check(err error) bool {
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func shasum(path string) []byte {
	fi, err := os.Open(path)
	if !check(err) {
		return nil
	}
	defer func() {
		if err := fi.Close(); err != nil {
			check(err)
		}
	}()
	hash := sha1.New()
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			if !check(err) {
				return nil
			}
		}
		if n == 0 {
			break
		}
		hash.Write(buf[:n])
	}
	return hash.Sum(nil)
}

func walk(target string, ignores []string) []string {
	results := make([]string, 0)
	limit := limiter.New()
	filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if !check(err) {
			return nil
		}
		// Ignore directory or file by pattern
		isDir := info.IsDir()
		if ignore(ignores, path) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}
		if isDir {
			return nil
		}
		// Calculate SHA1 of file
		limit.Add()
		go func() {
			defer limit.Done()
			sum := shasum(path)
			if sum != nil {
				result := fmt.Sprintf("%s, %x, %d\n", path, sum, info.Size())
				fmt.Print(result)
				results = append(results, result)
			}
		}()
		return nil
	})
	limit.Wait()
	sort.Strings(results)
	return results
}

func main() {
	// Define output path arguments
	output := flag.String("o", "", "")
	// Define command help
	flag.Usage = func() {
		fmt.Printf("%s\n\t%s\n",
			"Usage: dirwalker [path] [to/ignore/path]...",
			"-o Specific output file path, example: -o=result.txt")
	}
	flag.Parse()
	// Parse target directory arguments
	args := flag.Args()
	argsLen := len(args)
	target := filepath.Dir(os.Args[0])
	if argsLen > 0 {
		target = args[0]
	}
	absTarget, err := filepath.Abs(target)
	if !check(err) {
		return
	}
	// Parse ignored files / directories arguments
	ignores := make([]string, 0)
	if argsLen > 1 {
		for _, path := range args[1:] {
			path, err = filepath.Abs(path)
			if check(err) {
				ignores = append(ignores, path)
			}
		}
	}
	// Walk the specific directory
	results := walk(absTarget, ignores)
	if *output != "" {
		err = ioutil.WriteFile(*output, []byte(strings.Join(results, "")), 0644)
		check(err)
	}
}
