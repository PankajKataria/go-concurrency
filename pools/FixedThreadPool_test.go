package pools

import (
	"fmt"
	"testing"
	"time"
	u "github.com/PankajKataria/go-concurrency/utils"
)

func Test_HappyCaseResultShouldBeValid(t *testing.T) {
	ftp, _ := NewFixedThreadPool(1)
	f1, _ := ftp.Submit(func() interface{} {
		return 2.434
	})

	result, ok := f1.Result()
	if ok != true {
		t.Errorf("Expected true\nGot %v", ok)
	}

	if result != 2.434 {
		t.Errorf("Expected 2.434\nGot %v", result)
	}

	ftp.ShutDown()
	// Output:
	// Worker 0: Alive
	// Worker 0: Running Task 0
	// Worker 0: Finished Running Task 0
	// Worker 0: Shuting down
}

func Test_FaultyCallback(t *testing.T) {
	ftp, _ := NewFixedThreadPool(1)
	f1, _ := ftp.Submit(func() interface{} {
		// time.Sleep(2 *time.Second)
		return 2.434
	})

	f1.AddDoneCallback(func(*u.Future) {
		fmt.Println("Callback Called")
		var n2 int = 0
		var n1 int = 1 / n2
		fmt.Println(n1)
	})

	result, ok := f1.Result()
	if ok != true {
		t.Errorf("Expected true\nGot %v", ok)
	}

	if result != 2.434 {
		t.Errorf("Expected 2.434\nGot %v", result)
	}

	ftp.ShutDown()

	// Output:
	// Worker 0: Alive
	// Worker 0: Running Task 0
	// Worker 0: Finished Running Task 0
	// Worker 0: Shuting down
	// Callback Called
	// future 0 callback 0 errored
}

func Test_FaultyWork(t *testing.T) {
	ftp, _ := NewFixedThreadPool(1)
	f1, _ := ftp.Submit(func() interface{} {
		var n2 int = 0
		return 1 / n2
	})

	exception, ok := f1.Exception()
	if ok != true {
		t.Errorf("Expected true\nGot %v", ok)
	}

	expectedError := "runtime error: integer divide by zero"
	if exception != expectedError {
		t.Errorf("Expected [%s]\nGot [%v]", expectedError, exception)
	}

	ftp.ShutDown()

	// Output:
	// Worker 0: Alive
	// Worker 0: Running Task 0
	// Worker 0: Task  0 Paniced
	// Worker 0: Shuting down
}

func Test_FutureCancel(t *testing.T) {
	ftp, _ := NewFixedThreadPool(1)
	f1, _ := ftp.Submit(func() interface{} {
        time.Sleep(3 * time.Second)
        return 2.434
    })

    ok := f1.Cancel()
	if ok != true {
		t.Errorf("Expected true\nGot %v", ok)
	}

	ftp.ShutDown()

	// Output:
    // Worker 0: Alive
    // Worker 0: Task 0 is Cancelled
    // Worker 0: Shuting down
}
