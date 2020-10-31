package main
import (
	"fmt"
	"time"
	"github.com/PankajKataria/go-concurrency/pools"
	u "github.com/PankajKataria/go-concurrency/utils"

)

func main () {
	ftp := pools.NewFixedThreadPool(3)
	f1 := ftp.Submit (func () interface{} {
		time.Sleep(3 * time.Second)
		var n2 float32 = 2.434
		return n2
	})

	f1.AddDoneCallback(func (* u.Future) {
		fmt.Println("Callback Called")
	})

	err := f1.Exception()
	if (err != nil) {
		fmt.Printf("Exception Occured : %v\n", err)
	}

	result, ok := f1.Result()
	if (ok) {
		fmt.Printf("Result of task : %v\n", result)
	}

	ftp.ShutDown()
}