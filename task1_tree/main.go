package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
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

func dirTree(out io.Writer, path string, printFiles bool) error {
	if printFiles {
		return showDirF(out, path, "")
	}
	return showDirNoF(out, path, "")
}

func showDirF(out io.Writer, path string, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for i, file := range files {
		isLast := i == len(files)-1
		info, _ := file.Info()
		if isLast {
			if file.IsDir() {
				fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())
			} else {
				if info.Size() == 0 {
					fmt.Fprintf(out, "%s└───%s (empty)\n", prefix, file.Name())
				} else {
					fmt.Fprintf(out, "%s└───%s (%db)\n", prefix, file.Name(), info.Size())
				}
			}

		} else {
			if file.IsDir() {
				fmt.Fprintf(out, "%s├───%s\n", prefix, file.Name())
			} else {
				if info.Size() == 0 {
					fmt.Fprintf(out, "%s├───%s (empty)\n", prefix, file.Name())
				} else {
					fmt.Fprintf(out, "%s├───%s (%db)\n", prefix, file.Name(), info.Size())
				}
			}

		}

		if file.IsDir() {
			if isLast {
				err = showDirF(out, path+"/"+file.Name(), prefix+"\t")
			} else {
				err = showDirF(out, path+"/"+file.Name(), prefix+"│\t")
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func showDirNoF(out io.Writer, path string, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	filteredFiles := make([]fs.DirEntry, 0, len(files)/2)
	for _, val := range files {
		if val.IsDir() {
			filteredFiles = append(filteredFiles, val)
		}
	}

	for i, file := range filteredFiles {
		isLast := i == len(filteredFiles)-1
		if isLast {
			fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())

		} else {
			fmt.Fprintf(out, "%s├───%s\n", prefix, file.Name())
		}
		if isLast {
			err = showDirNoF(out, path+"/"+file.Name(), prefix+"\t")
		} else {
			err = showDirNoF(out, path+"/"+file.Name(), prefix+"│\t")
		}
		if err != nil {
			return err
		}

	}
	return nil
}
