package main

import (
	"fmt"
	"github.com/PankajKataria/go-concurrency/pools"
	u "github.com/PankajKataria/go-concurrency/utils"
)

func main() {
	ftp, err := pools.NewFixedThreadPool(10)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 100; i++ {
		f1, _ := ftp.Submit(func() interface{} {
			n2 := 2.434
			return n2
		})

		f1.AddDoneCallback(func(f *u.Future) {
			fmt.Printf("Callback Called %d\n", f.GetId())
			var n2 int = 0
			var n1 int = 1 / n2
			fmt.Println(n1)
		})

		result, ok := f1.Result()
		if ok {
			fmt.Printf("Result of task : %v\n", result)
		}
	}

	ftp.ShutDown()
}
