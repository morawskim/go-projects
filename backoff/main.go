package main

import (
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"time"
)

func main() {
	fmt.Println("Run exponential backoff without max retries, but with time limit")
	exponentialBackOff()

	fmt.Println("Run exponential backoff with max retries and time limit")
	exponentialBackOffWithMaxRetries()
}

func exponentialBackOff() {
	bo := createExponentialBackOff()

	err := backoff.Retry(func() error {
		fmt.Println(time.Now())
		return errors.New("some error")
	}, bo)

	if err != nil {
		fmt.Println(err)
	}
}

func exponentialBackOffWithMaxRetries() {
	bo := createExponentialBackOff()
	boWithMaxRetries := backoff.WithMaxRetries(bo, 3)

	err := backoff.Retry(func() error {
		fmt.Println(time.Now())
		return errors.New("some error")
	}, boWithMaxRetries)

	if err != nil {
		fmt.Println(err)
	}
}

func createExponentialBackOff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 20 * time.Second
	bo.Multiplier = 1.5
	bo.RandomizationFactor = 0.5

	return bo
}
