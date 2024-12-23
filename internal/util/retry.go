package util

import (
	"fmt"
	"time"
)

func RetryOperation(operation func() error, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		}
		time.Sleep(time.Duration(2^i) * time.Second)
	}
	return fmt.Errorf("operation failed after %d attempts", maxRetries)
}
