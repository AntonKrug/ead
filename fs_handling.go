package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func readFileEmbedded(filename string) (wholecontent string, size int64) {
	file, err := Assets.Open("/" + filename)
	checkErr(err)

	stats, err2 := file.Stat()
	checkErr(err2)
	defer file.Close()

	size = stats.Size()

	var filecontent []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		filecontent = append(filecontent, scanner.Text())
	}

	wholecontent = strings.Join(filecontent, NEW_LINE)
	return
}

func readFileReal(filename string) (wholecontent string, size int64) {
	content, err := ioutil.ReadFile(filename)
	checkErr(err)
	wholecontent = string(content)

	stats, err3 := os.Stat(filename)
	checkErr(err3)
	size = stats.Size()

	if *compressHtmlFlag && isCompressableExtension(filename) {
		var buf bytes.Buffer
		var gz *gzip.Writer
		gz, err = gzip.NewWriterLevel(&buf, gzip.BestCompression)

		_, err = gz.Write([]byte(wholecontent))
		checkErr(err)

		err = gz.Flush()
		checkErr(err)

		err = gz.Close()
		checkErr(err)

		log.Println("Compressed file", filename, "from", humanize.Bytes(uint64(size)), "to", humanize.Bytes(uint64(buf.Len())))

		wholecontent = string(buf.Bytes())
		size = int64(buf.Len())
	}

	return
}

func makePathRecursive(path string) {
	basepath := filepath.Dir(path)
	os.MkdirAll(basepath, os.ModePerm)
}

func writeContentToFile(path string, content string) {
	makePathRecursive(path)
	err := ioutil.WriteFile(path, []byte(content), 0664)
	checkErr(err)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func fileIsRealFile(filename string) bool {
	fileInfo, err := os.Stat(filename)
	checkErr(err)
	return !fileInfo.IsDir()
}

func isCompressableExtension(filename string) bool {
	supportedExtension := map[string]bool{
		"html": true,
		"htm":  true,
		"js":   true,
		"css":  true,
		"wasm": true,
		"xml":  true,
		"json": true,
	}

	ext := filepath.Ext(filename) // Get extension including the dot ".html"
	return supportedExtension[ext[1:]]
}
