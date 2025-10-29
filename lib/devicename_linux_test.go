//go:build linux

package lib

import (
	"sync"
	"testing"
)

var (
	testDevices = []Device{
		{Bus: 1, Path: []int{1}, VendorID: "1234", ProductID: "5678", Speed: "480M", DevNum: 1},
		{Bus: 1, Path: []int{1, 2}, VendorID: "5678", ProductID: "0001", Speed: "12M", DevNum: 2},
		{Bus: 2, Path: []int{1}, VendorID: "9999", ProductID: "8888", Speed: "480M", DevNum: 1},
	}

	testDevice = Device{
		Bus:       1,
		Path:      []int{1},
		VendorID:  "1234",
		ProductID: "5678",
		Speed:     "480M",
		DevNum:    1,
	}
)

// TestConcurrentEnrichAndClear tests concurrent access to deviceInfoCache
func TestConcurrentEnrichAndClear(t *testing.T) {
	const goroutines = 5
	const iterations = 20

	deviceInfoCache = make(map[string]deviceInfo) // Reset cache

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Readers
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				for _, dev := range testDevices {
					d := dev
					_ = d.enrich()
				}
			}
		}()
	}

	// Writers (clear cache)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				for _, dev := range testDevices {
					clearPriorityNameCache(dev)
				}
			}
		}()
	}

	wg.Wait()
}

// TestConcurrentEnumerateAndCache tests concurrent cache enumeration
func TestConcurrentEnumerateAndCache(t *testing.T) {
	const goroutines = 10

	deviceInfoCache = make(map[string]deviceInfo) // Reset cache

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			enumerateAndCache()
		}()
	}

	wg.Wait()
}

// TestConcurrentGetPriorityInfo tests concurrent reads of cache via getPriorityInfo
func TestConcurrentGetPriorityInfo(t *testing.T) {
	const goroutines = 10
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, _ = getPriorityInfo(testDevice)
			}
		}()
	}

	wg.Wait()
}

// TestConcurrentCacheClearAndRead tests clearing cache while reading
func TestConcurrentCacheClearAndRead(t *testing.T) {
	const readers = 10
	const clearers = 5
	const iterations = 50

	var wg sync.WaitGroup
	wg.Add(readers + clearers)

	enumerateAndCache()

	// Readers
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				for _, dev := range testDevices {
					_, _ = getPriorityInfo(dev)
				}
			}
		}()
	}

	// Clearers
	for i := 0; i < clearers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				for _, dev := range testDevices {
					clearPriorityNameCache(dev)
				}
			}
		}()
	}

	wg.Wait()
}
