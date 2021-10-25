package main

import (
	// "context"
	// "fmt"
	// "site/internal/http"
	// "site/internal/store/inmemory"
	"fmt"
	"site/test/compiler"
)

func main() {

	// store := inmemory.NewDB()

	// srv := http.NewServer(context.Background(), ":8080", store)

	// if err := srv.Run(); err != nil {
	// 	fmt.Println(err)
	// }

	// srv.WaitForGracefulTermination()
	err := compiler.BuildExe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
