package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	check(err)
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
			check(err)
		}
		if n == 0 {
			break
		}
		hash.Write(buf[:n])
	}
	return hash.Sum(nil)
}

func main() {
	// Parse command line arguments
	args := os.Args
	patterns := make([]string, 0)
	target := "./"
	argsLen := len(args)
	if argsLen > 1 {
		target = args[1]
	}
	if argsLen > 2 {
		for _, path := range args[2:] {
			patterns = append(patterns, path)
		}
	}
	absTarget, err := filepath.Abs(target)
	check(err)
	limiter := NewLimter()
	// Walk the specific directory
	filepath.Walk(absTarget, func(path string, info os.FileInfo, err error) error {
		check(err)
		// Ignore directory or file by pattern
		isDir := info.IsDir()
		if ignore(patterns, path) {
			if isDir {
				return filepath.SkipDir
			}
			return nil
		}
		if isDir {
			return nil
		}
		// Calculate SHA1 of file
		limiter.Add()
		go func() {
			sum := shasum(path)
			_, err = fmt.Printf("%s, %x, %d\n", path, sum, info.Size())
			check(err)
			limiter.Done()
		}()
		return nil
	})
	limiter.Wait()
}
