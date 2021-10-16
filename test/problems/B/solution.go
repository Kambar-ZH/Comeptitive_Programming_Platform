package main

import "fmt"


// Fibonachi problem
func main() {
	var N int
	fmt.Scanf("%d", &N)
	a, b := 1, 1
	for i := 0; i < N; i++ {
		a, b = b, a + b
	}
	fmt.Println(a)
}