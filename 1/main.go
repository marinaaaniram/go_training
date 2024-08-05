package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	err = recursiveTree(out, path, printFiles, "")
	return err
}

func recursiveTree(out io.Writer, path string, printFiles bool, indent string) (err error) {
	f, _ := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	files, _ := f.Readdir(-1)
	if err != nil {
		return err
	}

	files = sortFiles(files)

	if !printFiles {
		files = filterDirectories(files)
	}

	for i, file := range files {
		fileName := file.Name()
		sizeInfo := ""
		if !file.IsDir() {
			size := file.Size()
			if size > 0 {
				sizeInfo = " (" + strconv.FormatInt(size, 10) + "b)"
			} else {
				sizeInfo = " (empty)"
			}
		}

		if i < len(files)-1 {
			fmt.Fprintf(out, "%s├───%s%s\n", indent, fileName, sizeInfo)
		} else {
			fmt.Fprintf(out, "%s└───%s%s\n", indent, fileName, sizeInfo)
		}

		newIndent := indent
		if i < len(files)-1 {
			newIndent += "│\t"
		} else {
			newIndent += "\t"
		}
		path := path + "/" + fileName
		err := recursiveTree(out, path, printFiles, newIndent)
		if err != nil {
			return err
		}
	}
	return err
}

func sortFiles(files []os.FileInfo) []os.FileInfo {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	return files
}

func filterDirectories(files []os.FileInfo) []os.FileInfo {
	var dirs []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	return dirs
}
