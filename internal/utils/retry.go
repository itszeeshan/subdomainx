package utils

import (
	"fmt"
	"time"
)

func Retry[T any](fn func() (T, error), retries int, timeout int) (T, error) {
	var zero T
	var err error

	for i := 0; i < retries; i++ {
		result, fnErr := fn()
		if fnErr == nil {
			return result, nil
		}
		err = fnErr

		// Exponential backoff
		backoff := time.Duration(i*i) * time.Second
		if backoff > time.Duration(timeout)*time.Second {
			backoff = time.Duration(timeout) * time.Second
		}

		time.Sleep(backoff)
	}

	return zero, fmt.Errorf("after %d retries: %v", retries, err)
}
