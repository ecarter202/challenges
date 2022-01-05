package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("no filepath given")
	}

	start := time.Now()

	input := args[1]

	f, err := os.Open(input)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	defer f.Close()

	r := compress(f)

	output, err := os.Create(outputFile(args))
	if err != nil {
		log.Fatalf("error creating output file: %v", err)
	}

	if _, err := io.Copy(output, r); err != nil {
		log.Fatalf("error writing gzip: %v", err)
	}

	fmt.Printf("Compressed file in %dms\n", time.Since(start).Milliseconds())
}

func compress(rc io.ReadCloser) io.Reader {
	defer rc.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, rc); err != nil {
		log.Fatalf("error writing to buffer: %v", err)
	}

	outBuffer := new(bytes.Buffer)
	writer, err := gzip.NewWriterLevel(outBuffer, gzip.BestSpeed)
	if err != nil {
		log.Fatalf("error creating gzip writer: %v", err)
	}
	defer writer.Close()

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		log.Fatalf("error writing to buffer: %v", err)
	}

	return bytes.NewReader(outBuffer.Bytes())
}

func outputFile(args []string) (output string) {
	if len(args) < 2 {
		return ""
	} else if len(args) > 2 {
		return fmt.Sprintf("%s.gzip", args[2])
	}

	inputX := strings.Split(args[1], ".")

	if len(inputX) > 1 { // has ext
		return fmt.Sprintf("%s_out.%s.gzip", inputX[0], inputX[1])
	}

	return fmt.Sprintf("%s_out.gzip", inputX[0])
}
