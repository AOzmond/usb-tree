//go:build linux

package lib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jochenvg/go-udev"
)

type deviceInfo struct {
	Name  string
	Speed string
}

var (
	deviceInfoCache   = map[string]deviceInfo{}
	deviceInfoCacheMu sync.RWMutex
)

func (d *Device) enrich() bool {
	info, ok := getPriorityInfo(*d)
	if !ok {
		return false
	}

	if len(strings.TrimSpace(info.Name)) > 0 {
		d.Name = info.Name
	}
	d.Speed = info.Speed
	return true
}

func (d *Device) getPriorityNameCacheKey() string {
	return fmt.Sprintf("%s:%s:%03d:%03d", d.VendorID, d.ProductID, d.Bus, d.DevNum)
}

func getPriorityInfo(device Device) (deviceInfo, bool) {
	key := device.getPriorityNameCacheKey()

	deviceInfoCacheMu.RLock()
	info, found := deviceInfoCache[key]
	deviceInfoCacheMu.RUnlock()

	if !found {
		enumerateAndCache()
	}

	if found {
		return info, true
	}

	return deviceInfo{}, false
}

func enumerateAndCache() {
	u := udev.Udev{}
	e := u.NewEnumerate()
	err := e.AddMatchSubsystem("usb")
	if err != nil {
		return
	}

	devices, _ := e.Devices()

	newCache := make(map[string]deviceInfo)

	for _, device := range devices {
		vid := device.PropertyValue("ID_VENDOR_ID")

		if vid != "" {
			vendorName := device.PropertyValue("ID_VENDOR_FROM_DATABASE")
			if vendorName == "" {
				vendorName = device.SysattrValue("manufacturer")
			}

			deviceName := device.PropertyValue("ID_MODEL_FROM_DATABASE")
			if deviceName == "" {
				deviceName = device.SysattrValue("product")
			}

			name := vendorName + " " + deviceName
			pid := device.PropertyValue("ID_MODEL_ID")
			bus := device.PropertyValue("BUSNUM")
			devNum := device.PropertyValue("DEVNUM")
			speed := device.SysattrValue("speed")

			key := fmt.Sprintf("%s:%s:%03s:%03s", vid, pid, bus, devNum)

			newCache[key] = deviceInfo{
				Name:  name,
				Speed: speed,
			}
		}
	}

	// Replace the entire cache with write lock (minimize lock time)
	deviceInfoCacheMu.Lock()
	deviceInfoCache = newCache
	deviceInfoCacheMu.Unlock()
}

func clearPriorityNameCache(device Device) {
	key := device.getPriorityNameCacheKey()

	deviceInfoCacheMu.Lock()
	delete(deviceInfoCache, key)
	deviceInfoCacheMu.Unlock()
}
