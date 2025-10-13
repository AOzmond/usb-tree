//go:build !linux

package lib

func (d *Device) enrich() bool {
	return true
}

func clearPriorityNameCache(device Device) {
	return
}
