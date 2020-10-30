package main
import (
	"fmt"
	"time"
	"github.com/PankajKataria/go-concurrency/pools"
	f "github.com/PankajKataria/go-concurrency/future"

)

func main () {
	ftp := pools.NewFixedThreadPool(3)
	f1 := ftp.Submit (func () interface{} {
		time.Sleep(3 * time.Second)
		return "Hello"
	})

	f1.AddDoneCallback(func (* f.Future) {
		fmt.Println("callback Called")
	})

	f1.Running()
	ftp.ShutDown()
}