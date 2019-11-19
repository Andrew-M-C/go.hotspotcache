package hotspotcache

import (
	"log"
	"runtime"
	"sync"
	"testing"
)

func init() {
	log.Printf("runtime.GOMAXPROCS(0) = %d\n", runtime.GOMAXPROCS(0))
}

func BenchmarkMultiRoutine(b *testing.B) {
	c := New(b.N / 4)
	mask := 1
	for {
		mask <<= 1
		if mask >= b.N/8 {
			break
		}
	}

	MAX := runtime.GOMAXPROCS(0) * 2
	var wg sync.WaitGroup
	wg.Add(MAX)
	b.ResetTimer()

	for i := 0; i < MAX; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < b.N; i++ {
				n := i & mask
				_, exist := c.Load(n)
				if false == exist {
					c.Store(n, i)
				}
			}
		}()
	}

	wg.Wait()
	b.Logf("b.N = %d", b.N)
	return
}

func Benchmark75PercentWrite(b *testing.B) {
	c := New(b.N / 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n := i
		_, exist := c.Load(n)
		if false == exist {
			c.Store(n, i)
		}
	}

	b.Logf("b.N = %d", b.N)
	return
}

func BenchmarkNoHitNoWrite(b *testing.B) {
	c := New(b.N / 4)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n := i
		_, _ = c.Load(n)
	}

	b.Logf("b.N = %d", b.N)
	return
}

func BenchmarkAllHit(b *testing.B) {
	c := New(b.N / 4)
	c.Store(1, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = c.Load(1)
	}

	b.Logf("b.N = %d", b.N)
	return
}
