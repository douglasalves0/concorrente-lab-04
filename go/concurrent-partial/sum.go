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

type MyADT2 struct {
	path1      string
	path2      string
	similarity float64
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

func Min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func similarity(s1 string, s2 string, c chan MyADT2) {
	data1, _ := readFile(s1)
	data2, _ := readFile(s2)
	sums1 := make(map[int]int)
	sums2 := make(map[int]int)

	var total float64
	total = 0
	for _, b := range data1 {
		sums1[int(b)]++
		total++
	}
	for _, b := range data2 {
		sums2[int(b)]++
	}
	var counter float64
	counter = 0

	for k, v := range sums1 {
		counter += float64(Min(v, sums2[k]))
	}

	ans := MyADT2{s1, s2, counter / total}
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

	c2 := make(chan MyADT2)
	for i, path1 := range os.Args[1:] {
		for _, path2 := range os.Args[i+2:] {
			go similarity(path1, path2, c2)
		}
	}

	for i := range os.Args[1:] {
		for range os.Args[i+2:] {
			ans := <-c2
			fmt.Printf("Similarity between %s and %s: %.5f", ans.path1, ans.path2, ans.similarity)
			fmt.Print("%\n")
		}
	}
}
