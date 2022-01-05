package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	dir, comma, codeColumn string
	hasHeader, getAll      bool
	codeColIndex           int
)

func readDir(path, extType string) (files []os.FileInfo, err error) {
	files, err = ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	if extType != "" {
		for i, f := range files {
			x := strings.Split(f.Name(), ".")
			if x[len(x)-1] != extType {
				files = append(files[:i], files[i+1:]...)
			}
		}
	}

	return files, nil
}

func fileReader(ctx context.Context, file *os.File, codeChan chan<- string) {
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = []rune(comma)[0]

	var iRow int
	for {
		if iRow == 0 && hasHeader {
			csvReader.Read() // discard header value
			iRow++
			continue
		}

		select {
		case <-ctx.Done():
			break
		default:
		}

		rec, err := csvReader.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error reading file: %v\n", err)
			}
			break
		}

		code := rec[codeColIndex]

		codeChan <- code

		iRow++
	}
}

func duplicateChecker(codeChan <-chan string, duplicates *[]string, cancel context.CancelFunc) {
	seen := new(sync.Map)

	for code := range codeChan {
		if _, exists := seen.LoadOrStore(code, true); exists {
			*duplicates = append(*duplicates, code)
			if !getAll {
				cancel()
			}
		}
	}
}

func main() {
	flag.StringVar(&dir, "d", "", "Specifies the directory containing the CSV files for processing.")
	flag.StringVar(&comma, "s", ",", "Specifies the separator \"comma\" used in the files.")
	flag.StringVar(&codeColumn, "c", "A", "Which column contains the code values.")
	flag.BoolVar(&hasHeader, "h", true, "Indicates if the files contain header rows which should be skipped.")
	flag.BoolVar(&getAll, "a", false, "Get all duplicates instead of killing after finding the first.")
	flag.Parse()

	codeColIndex = columnLetterToIndex(codeColumn)

	start := time.Now()

	files, err := readDir(dir, "csv")
	if err != nil {
		log.Fatalf("reading directory: %v", err)
	}

	var openedFiles []*os.File
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
			continue
		}

		openedFiles = append(openedFiles, f)
	}

	var duplicates []string
	if len(openedFiles) > 0 {
		codeChan := make(chan string, 100)
		wg := new(sync.WaitGroup)
		wg2 := new(sync.WaitGroup)
		ctx, cancel := context.WithCancel(context.Background())

		wg2.Add(1)
		go func() {
			defer wg2.Done()
			duplicateChecker(codeChan, &duplicates, cancel)
		}()

		for i := 0; i < len(openedFiles); i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				fileReader(ctx, openedFiles[index], codeChan)
			}(i)
		}

		wg.Wait()
		close(codeChan)
		wg2.Wait()
	}

	fmt.Println()
	fmt.Println("##################################")
	if getAll {
		fmt.Printf("# Duplicates found: %d\n", len(duplicates))
		fmt.Printf("# %v\n", duplicates)
	} else {
		fmt.Printf("# Duplicate found: %v\n", duplicates[0])
	}
	fmt.Printf("# Runtime: %vms\n", time.Since(start).Milliseconds())
	fmt.Println("##################################")
	fmt.Println()
}
