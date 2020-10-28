# go-concurrency
Fun stuff implemented in go

## Usage
- Example
```go
	ftp := pools.NewFixedThreadPool(5) // Creating thread pool of 5 workers
    	f1 := ftp.Submit (func () interface{} { // Submiting an anonymous function
		time.Sleep(10 * time.Second)
		return "Hello"
	})
```

- Functions supported by fixed thread pool 
``` go
	NewFixedThreadPool(workers int) (*FixedThreadPool)
	Submit (callable func() (interface{})) *f.Future 
	ShutDown () // Graceful shutdown
```

- Functions supported by Future
``` go
	f1.Cancel () bool   // Cancel the submitted task
	f1.Cancelled () bool // Check if the task is cancelled
	f1.Running () bool 
	f1.Finished () bool
	f1.Result(timeout int) (interface{}, bool) // Get the result if not available wait timeout seconds
```

- Supported Features
    - [X] [Fixed thread pools](pools/fixedThreadPool.go)
    - [X] [Futures](future/future.go)

## Todo
- [ ] Scheduled thread pool

## Good Links
 - [Using Interfaces in Go](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)
