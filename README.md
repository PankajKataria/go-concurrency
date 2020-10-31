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
	ftp.ShutDown()
```

- Functions supported by fixed thread pool 
``` go
	NewFixedThreadPool(workers int) ((*FixedThreadPool), error)
	Submit (callable func() (interface{})) (*f.Future, error)
	ShutDown () // Graceful shutdown
```

- Functions supported by Future
``` go
	// Check if the task status
	f1.Cancel () bool
	f1.Cancelled () bool
	f1.Running () bool 
	f1.Finished () bool
	// Get the result if not available wait timeout seconds
	f1.Result(timeout ...int) (interface{}, bool) 
	// Get the exception if any occured
	f1.Exception(timeout ...int) (interface{}, bool)
	// Add Callback to be executed after task is Finised/Canceled
	f1.AddDoneCallback(callback func(*Future))
```

- Supported Features
    - [X] [Fixed thread pools](pools/fixedThreadPool.go)
    - [X] [Futures](utils/future.go)

## Todo
- [ ] Scheduled thread pool

## Good Links
 - [Using Interfaces in Go](https://jordanorelli.com/post/32665860244/how-to-use-interfaces-in-go)
