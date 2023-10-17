// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdownHandler(t *testing.T) {
	shutdown := NewShutdownHandler()
	var wg sync.WaitGroup
	wg.Add(2)
	count := 0
	shutdown.AddListener(func() {
		count += 1
		wg.Done()
	})

	exit := func() {
		count += 1
		wg.Done()
	}

	shutdown.Listen(exit)
	shutdown.Interrupt()

	wg.Wait()
	assert.Equal(t, 2, count)
}
