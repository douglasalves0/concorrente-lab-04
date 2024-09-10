package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

type MyADT struct {
	sum  int
	path string
}

// sum all bytes of a file
func sum(filePath string, c chan MyADT) {
	data, _ := readFile(filePath)

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	ans := MyADT{_sum, filePath}
	c <- ans
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	sums := make(map[int][]string)
	c := make(chan MyADT)
	for _, path := range os.Args[1:] {
		go sum(path, c)
	}
	for range os.Args[1:] {
		_sum := <-c
		totalSum += int64(_sum.sum)
		sums[_sum.sum] = append(sums[_sum.sum], _sum.path)
	}

	fmt.Println(totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
