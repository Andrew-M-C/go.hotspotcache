package hotspotcache

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestBasic(t *testing.T) {
	const N = 20000
	const SIZE = 0 // causes default as 10240
	c := New(SIZE)

	for i := 0; i < SIZE; i++ {
		c.Store(i, i)
	}

	time.Sleep(time.Second)

	for i := 0; i < SIZE; i++ {
		res, exist := c.Load(i)
		if false == exist {
			t.Errorf("Key %v not exist", i)
			return
		}
		if res.(int) != i {
			t.Errorf("Key %v not equal to %v", res, i)
		}
	}

	nonexistCount := 0

	for i := 0; i < N; i++ {
		n := int(rand.Int31n(N))
		_, exist := c.Load(n)
		if false == exist {
			nonexistCount++
			c.Store(n, n)
		}
	}

	t.Logf("max size: %d", c.MaxSize())
	t.Logf("nonexist percentage: %.02f%%", float64(nonexistCount)/float64(N)*100)

	return
}

func TestConcurrency(t *testing.T) {
	const N = 1000000
	const SIZE = 50000
	c := New(SIZE)

	MAX := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	wg.Add(MAX)

	for i := 0; i < MAX; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < N; i++ {
				n := int(rand.Int31n(N))
				_, exist := c.Load(n)
				if false == exist {
					c.Store(n, n)
				}
			}
		}()
	}

	time.Sleep(10 * time.Millisecond)
	t.Logf("inter status: %s", c.dumpStatus())

	wg.Wait()
	t.Logf("status: %s", c.dumpStatus())
	time.Sleep(time.Second)
	t.Logf("status: %s", c.dumpStatus())
	return
}
