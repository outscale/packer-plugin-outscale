package retry

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	numTries := uint(0)
	// Test that a passing function only gets called once.
	err := Run(0, 0, 0, func(i uint) (bool, error) {
		numTries++
		return true, nil
	})
	if numTries != 1 {
		t.Fatal("Passing function should not have been retried.")
	}
	if err != nil {
		t.Fatalf("Passing function should not have returned a retry error. Error: %s", err)
	}

	// Test that a failing function gets retried (once in this example).
	numTries = 0
	results := []bool{false, true}
	err = Run(0, 0, 0, func(i uint) (bool, error) {
		result := results[numTries]
		numTries++
		return result, nil
	})
	if numTries != 2 {
		t.Fatalf("Retried function should have been tried twice. Tried %d times.", numTries)
	}
	if err != nil {
		t.Fatalf("Successful retried function should not have returned a retry error. Error: %s", err.Error())
	}

	// Test that a function error gets returned, and the function does not get called again.
	numTries = 0
	funcErr := errors.New("this function had an error!")
	err = Run(0, 0, 0, func(i uint) (bool, error) {
		numTries++
		return false, funcErr
	})
	if numTries != 1 {
		t.Fatal("Errant function should not have been retried.")
	}
	if err != funcErr {
		t.Fatalf("Errant function did not return the right error %s. Error: %s", funcErr, err.Error())
	}

	// Test when a function exhausts its retries.
	numTries = 0
	expectedTries := uint(3)
	err = Run(0, 0, expectedTries, func(i uint) (bool, error) {
		numTries++
		return false, nil
	})
	if numTries != expectedTries {
		t.Fatalf("Unsuccessful retry function should have been called %d times. Only called %d times.", expectedTries, numTries)
	}
	if err != ErrRetryExhausted {
		t.Fatalf("Unsuccessful retry function should have returned a retry exhausted error. Actual error: %s", err.Error())
	}
}
