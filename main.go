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

func main() {
	// Parse target directory arguments
	args := os.Args
	argsLen := len(args)
	target := filepath.Dir(args[0])
	if argsLen > 1 {
		target = args[1]
	}
	absTarget, err := filepath.Abs(target)
	if !check(err) {
		return
	}
	// Parse ignored directories arguments
	patterns := make([]string, 0)
	if argsLen > 2 {
		for _, path := range args[2:] {
			patterns = append(patterns, path)
		}
	}
	// Walk the specific directory
	limiter := NewLimter()
	filepath.Walk(absTarget, func(path string, info os.FileInfo, err error) error {
		if !check(err) {
			return nil
		}
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
			defer limiter.Done()
			sum := shasum(path)
			if sum != nil {
				_, err = fmt.Printf("%s, %x, %d\n", path, sum, info.Size())
				check(err)
			}
		}()
		return nil
	})
	limiter.Wait()
}
