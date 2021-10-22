package main

import "fmt"


// Fibonachi problem
func main() {
	var N int
	fmt.Scanf("%d", &N)
	F := make([]int, N + 1)
	F[0] = 1
	F[1] = 1
	for i := 2; i <= N; i++ {
		F[i] = F[i - 1] + F[i - 2]
	}
	fmt.Println(F[N])
}