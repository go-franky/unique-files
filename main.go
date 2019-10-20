package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type arrayDirectory []string

func (d *arrayDirectory) String() string {
	return "An list of directories to look for files"
}

func (d *arrayDirectory) Set(value string) error {
	*d = append(*d, value)
	return nil
}

// main runs the program to list all file in a given directory
// whose contents are the same.
// Example:
//    go run main.go --path path/to/folder/1 --path path/to/folder/2
func main() {
	var dirs arrayDirectory
	flag.Var(&dirs, "path", "path to find duplicate files, defaults to current directory")
	flag.Parse()
	if len(dirs) == 0 {
		dirs = append(dirs, ".")
	}

	filePaths := make(chan string)
	go getFiles(filePaths, dirs...)

	files := map[string][]string{}

	for path := range filePaths {
		md5, err := hashMD5File((path))
		if err != nil {
			log.Printf("could not create MD5 of file %s: %v", path, err)
		}
		if p, ok := files[md5]; ok {
			files[md5] = append(p, path)
		} else {
			files[md5] = []string{path}
		}
	}

	for md5, files := range files {
		fmt.Printf("MD5: %s\n\t%s\n", md5, strings.Join(files, "\n\t"))
	}
}

func hashMD5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "could not open file")
	}
	defer file.Close()
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return "", errors.Wrap(err, "could not copy file")
	}
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}

// getFiles thakes a root directory, and sends all files
// that are not directories to the channel
func getFiles(paths chan<- string, dirs ...string) {
	defer close(paths)
	for _, root := range dirs {

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			paths <- path
			return nil
		})
		if err != nil {
			log.Printf("could not read file: %v", err)
		}
	}
}
