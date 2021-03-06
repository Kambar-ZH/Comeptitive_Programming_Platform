package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader *bufio.Reader = bufio.NewReader(os.Stdin)

func scanf(f string, a ...interface{}) { fmt.Fscanf(reader, f, a...) }

// Sum on the array problem
func main() {
	var N int
	fmt.Scanf("%d\n", &N)
	Arr := make([]int, N)
	var i int = 0
	for ; i < N; i++ {
		scanf("%d", &Arr[i])
	}
	var sum int
	for _, x := range Arr {
		sum += x
	}
	fmt.Println(sum)
}
