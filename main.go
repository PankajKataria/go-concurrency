package main
import (
	"fmt"
	"github.com/PankajKataria/go-concurrency/pools"
)

func main () {
	ftp := pools.NewFixedThreadPool(5)
	for i := 0; i<10; i++ {
		f1 := ftp.Submit (func () []interface{} {
			fmt.Println("Hello World")
			return "Hello", "World"
		})
		fmt.Println(f1.Cancelled())
	}
}