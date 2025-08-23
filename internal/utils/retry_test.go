package utils

import (
	"errors"
	"testing"
	"time"
)

func TestRetrySuccessOnFirstTry(t *testing.T) {
	called := false
	fn := func() (string, error) {
		called = true
		return "success", nil
	}

	result, err := Retry(fn, 3, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "success" {
		t.Errorf("Expected 'success', got %s", result)
	}

	if !called {
		t.Error("Function should have been called")
	}
}

func TestRetrySuccessOnSecondTry(t *testing.T) {
	callCount := 0
	fn := func() (int, error) {
		callCount++
		if callCount == 1 {
			return 0, errors.New("first attempt failed")
		}
		return 42, nil
	}

	result, err := Retry(fn, 3, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryAllAttemptsFail(t *testing.T) {
	callCount := 0
	expectedError := errors.New("persistent error")
	fn := func() (bool, error) {
		callCount++
		return false, expectedError
	}

	result, err := Retry(fn, 3, 10)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != false {
		t.Errorf("Expected false, got %v", result)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", callCount)
	}

	// Check error message
	expectedMsg := "after 3 retries: persistent error"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestRetryZeroRetries(t *testing.T) {
	called := false
	fn := func() (string, error) {
		called = true
		return "", errors.New("should not be called")
	}

	result, err := Retry(fn, 0, 10)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}

	if called {
		t.Error("Function should not be called with 0 retries")
	}
}

func TestRetryWithStruct(t *testing.T) {
	type TestStruct struct {
		Name  string
		Value int
	}

	callCount := 0
	fn := func() (TestStruct, error) {
		callCount++
		if callCount == 1 {
			return TestStruct{}, errors.New("first attempt failed")
		}
		return TestStruct{Name: "test", Value: 123}, nil
	}

	result, err := Retry(fn, 3, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := TestStruct{Name: "test", Value: 123}
	if result != expected {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryWithSlice(t *testing.T) {
	callCount := 0
	fn := func() ([]string, error) {
		callCount++
		if callCount == 1 {
			return nil, errors.New("first attempt failed")
		}
		return []string{"a", "b", "c"}, nil
	}

	result, err := Retry(fn, 3, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := []string{"a", "b", "c"}
	if len(result) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(result))
	}

	for i, item := range expected {
		if result[i] != item {
			t.Errorf("Expected item %d to be %s, got %s", i, item, result[i])
		}
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryBackoffTiming(t *testing.T) {
	start := time.Now()
	callCount := 0
	fn := func() (int, error) {
		callCount++
		return 0, errors.New("always fail")
	}

	_, err := Retry(fn, 3, 5)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	elapsed := time.Since(start)

	// With 3 retries and exponential backoff (0s, 1s, 4s), minimum time should be ~5 seconds
	// But we cap at timeout (5s), so it should be around 5 seconds
	if elapsed < 4*time.Second {
		t.Errorf("Expected at least 4 seconds of backoff, got %v", elapsed)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", callCount)
	}
}
