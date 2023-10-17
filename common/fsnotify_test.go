// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFsNotify(t *testing.T) {
	watcher := NewFileWatcher()
	var wg sync.WaitGroup
	wg.Add(1)
	count := 0
	os.Create("/tmp/1")
	watcher.Add("/tmp/1", func() {

		count += 1
		wg.Done()
	})
	watcher.Listen()
	os.Remove("/tmp/1")
	wg.Wait()
	assert.Equal(t, 1, count)

	// watched file DNE
	os.Remove("/tmp/1")
	err := watcher.Add("/tmp/1", func() {})
	assert.NotNil(t, err)
}
