package main
import (
	"time"
	"github.com/PankajKataria/go-concurrency/pools"
)

func main () {
	ftp := pools.NewFixedThreadPool(5)
	f1 := ftp.Submit (func () interface{} {
		time.Sleep(10 * time.Second)
		return "Hello"
	})
	f1.Running()
	ftp.ShutDown()
}