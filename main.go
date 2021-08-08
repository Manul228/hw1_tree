package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func dirTree(out io.Writer, path string, printFiles bool) error {

	tree, err := getDirContent(path, 0, 0, printFiles)
	fmt.Fprint(out, tree)

	return err
}

func getDirContent(dirPath string, depth int, stopDepth int, printFiles bool) (string, error) {

	var result string
	var part string

	files, err := ioutil.ReadDir(dirPath)
	listLength := len(files)

	for idx, file := range files {
		if !file.IsDir() && !printFiles {
			continue
		}
		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			result += setOffsets(file, idx, listLength, depth, stopDepth)
			if idx == listLength-1 {
				stopDepth++
			}
			part, _ = getDirContent(filePath, depth+1, stopDepth, printFiles)
			result += part
		} else {
			result += setOffsets(file, idx, listLength, depth, stopDepth)
		}
	}
	return result, err
}

func setOffsets(file fs.FileInfo, idx int, listLength int, depth int, stopDepth int) string {
	var offset string
	var result string

	offset += strings.Repeat("│\t", depth-stopDepth)
	offset += strings.Repeat("\t", stopDepth)

	info := getFileStatsStr(file)

	if idx == listLength-1 {
		result = offset + "└───"
	} else {
		result = offset + "├───"
	}

	return result + info
}

func getFileStatsStr(file fs.FileInfo) string {
	var size string

	if !file.IsDir() {
		if file.Size() == 0 {
			size = "empty"
		} else {
			size = strconv.FormatInt(file.Size(), 10) + "b"
		}
		return fmt.Sprintf("%v (%v)\n", file.Name(), size)
	}
	return fmt.Sprintf("%v\n", file.Name())
}

func main() {
	// out := os.Stdout
	// if !(len(os.Args) == 2 || len(os.Args) == 3) {
	// 	panic("usage go run main.go . [-f]")
	// }
	// path := os.Args[1]
	// printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// err := dirTree(out, path, printFiles)
	// if err != nil {
	// 	panic(err.Error())
	// }

	tree, _ := getDirContent("testdata", 0, 0, false)
	fmt.Println(tree)
}
