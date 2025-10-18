//go:build linux

package lib

import (
	"fmt"
	"strings"

	"github.com/jochenvg/go-udev"
)

type deviceInfo struct {
	Name  string
	Speed string
}

var deviceInfoCache = map[string]deviceInfo{}

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

	if _, found := deviceInfoCache[key]; !found {
		enumerateAndCache()
	}

	if info, found := deviceInfoCache[key]; found {
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

			deviceInfoCache[key] = deviceInfo{
				Name:  name,
				Speed: speed,
			}
		}
	}
}

func clearPriorityNameCache(device Device) {
	key := device.getPriorityNameCacheKey()
	delete(deviceInfoCache, key)
}
