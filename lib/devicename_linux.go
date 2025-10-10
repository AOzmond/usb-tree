package lib

import (
	"fmt"
	"time"

	"github.com/jochenvg/go-udev"
)

type deviceInfo struct {
	Name  string
	Speed string
}

var deviceInfoCache = map[string]deviceInfo{}

func (d *Device) getPriorityNameCacheKey() deviceInfo {
	// Check cache for device name

	key := fmt.Sprintf("%s:%s:%03d:%03d", d.VendorID, d.ProductID, d.Bus, d.DevNum)
	if info, found := deviceInfoCache[key]; found {
		return info
	}

	// If not found, get from udev
	enumerateAndCache()

	if info, found := deviceInfoCache[key]; found {
		return info
	}

	addErrorLog(">"+key+"< not found in cache", time.Now(), StateError)
	// Return empty deviceInfo if still not found
	return deviceInfo{}
}

func enumerateAndCache() {
	u := udev.Udev{}
	e := u.NewEnumerate()
	err := e.AddMatchSubsystem("usb")
	if err != nil {
		addErrorLog(err.Error(), time.Now(), StateError)
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
			devnum := device.PropertyValue("DEVNUM")
			speed := device.SysattrValue("speed")

			key := fmt.Sprintf("%s:%s:%03s:%03s", vid, pid, bus, devnum)
			deviceInfoCache[key] = deviceInfo{
				Name:  name,
				Speed: speed,
			}
		}
	}
}
