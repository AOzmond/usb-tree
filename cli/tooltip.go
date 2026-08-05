package cli

import "strconv"

// getSelectedDeviceInfo returns formatted device info for the currently selected node
func (m *Model) getSelectedDeviceInfo() string {
	if m.selectedDevice == nil {
		return ""
	}
	node := m.selectedDevice

	busStyle := windowStyle.Foreground(busTextColor)
	deviceStyle := windowStyle.Foreground(deviceTextColor)
	vidStyle := windowStyle.Foreground(vidTextColor)
	pidStyle := windowStyle.Foreground(pidTextColor)
	nameStyle := windowStyle.Foreground(nameTextColor)
	linkStyle := windowStyle.Foreground(linkTextColor)

	busString := busStyle.Render("Bus: ", strconv.Itoa(node.Bus))
	deviceString := deviceStyle.Render(" Device: ", strconv.Itoa(node.DevNum))
	vidString := vidStyle.Render(" VID: ", node.VendorID)
	pidString := pidStyle.Render(" PID: ", node.ProductID)

	deviceInfo := busString + deviceString + vidString + pidString

	nameString := nameStyle.Render(node.Name)
	linkString := linkStyle.Render(getDbAddress(node.VendorID, node.ProductID))

	tooltipString := deviceInfo + "\n" + nameString + "\n" + linkString

	return tooltipString
}

// getDbAddress returns the USB-ID database link for the given VID and PID
func getDbAddress(vid string, pid string) string {
	baseAddress := "https://the-sz.com/products/usbid/?v="
	return baseAddress + vid + "&p=" + pid
}
