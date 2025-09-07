package usbtreelib

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
)

// A LogState represents the difference between the cached Device and current Device.
type LogState string

// A Device represents a USB Device
type Device struct {
	Path      []int
	Name      string
	VendorID  string
	ProductID string
	Speed     string
	Bus       int
	State     LogState
}

// TreeNode represents a Device and its children for building tree structures.
type TreeNode struct {
	Device
	Children []TreeNode
}

// Log represents a change in a Device.
type Log struct {
	Time  time.Time
	Text  string
	State LogState
}

// These constants represent the State of a Device.
const (
	StateNormal  LogState = "normal"
	StateAdded   LogState = "added"
	StateRemoved LogState = "removed"
)

var (
	cachedDevices []Device
	lastMergedMap map[string]Device
	logs          []Log
)

var pollingStop chan bool

// Stop will turn off the polling of new Devices.
func Stop() {
	select {
	case pollingStop <- true:
	default:
	}
}

// Init will start the polling of Devices connected to the machine, and return the current list of connected Devices.
// It takes a callback function to be run anytime there is a change in Devices
func Init(onUpdateCallback func([]Device)) []Device {
	pollingStop = make(chan bool)
	Refresh()
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logtime, newDevices := getDevices()
				changed, mergedDevices := deviceDiff(newDevices, logtime)
				if changed {
					onUpdateCallback(mergedDevices)
				}

			case <-pollingStop:
				close(pollingStop)
				return
			}
		}
	}()

	return cachedDevices
}

// Refresh resets the cached Device state to that of the current devices connected to the machine.
func Refresh() []Device {
	_, cachedDevices = getDevices()
	return cachedDevices
}

// returns lists of devices in depth first search order.
func getDevices() (time.Time, []Device) {
	ctx := gousb.NewContext()
	defer ctx.Close()
	devices := []Device{}

	_, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		devices = append([]Device{descToDevice(*desc)}, devices...)
		return false
	})
	if err != nil {
		log.Printf("Issue accessing USB Devices: %v", err)
		return time.Now(), nil
	}
	return time.Now(), devices
}

// Returns device based on a given DeviceDesc
func descToDevice(desc gousb.DeviceDesc) Device {
	return Device{
		Bus:       desc.Bus,
		Path:      desc.Path,
		Name:      usbid.Describe(&desc),
		VendorID:  desc.Vendor.String(),
		ProductID: desc.Product.String(),
		Speed:     desc.Speed.String(),
		State:     StateNormal,
	}
}

func (d *Device) deviceKey() string {
	return fmt.Sprintf("%d:%v:%s:%s:%s", d.Bus, d.Path, d.VendorID, d.ProductID, d.Speed)
}

func deviceDiff(newDevices []Device, logtime time.Time) (changed bool, merged []Device) {
	mergedMap := make(map[string]Device)
	changed = false

	// Mark all cachedDevices as removed initially
	for _, device := range cachedDevices {
		device.State = StateRemoved
		mergedMap[device.deviceKey()] = device
	}

	// Reset persisting devices to normal, and add new devices
	for _, newDevice := range newDevices {
		key := newDevice.deviceKey()
		if existingDevice, exists := mergedMap[key]; exists {
			// Device exists, reset its status to normal
			existingDevice.State = StateNormal
			mergedMap[key] = existingDevice
		} else {
			// Device is new, add to mergedMap
			newDevice.State = StateAdded
			mergedMap[key] = newDevice
		}
	}

	// Convert map back into slice and log changes since lastMergedMap
	merged = make([]Device, 0, len(mergedMap))
	for key, device := range mergedMap {
		merged = append(merged, device)

		if lastDevice, exists := lastMergedMap[key]; !exists {
			addDeviceLog(device, logtime)
			changed = true
		} else if device.State != lastDevice.State {
			addDeviceLog(device, logtime)
			changed = true
		}
	}

	merged = sortDevices(merged)

	lastMergedMap = mergedMap

	return changed, merged
}

// BuildDeviceTree takes a []Device and converts each Device to a TreeNode and returns the roots of those Trees.
func BuildDeviceTree(devices []Device) []TreeNode {
	roots := []TreeNode{}
	nodes := []TreeNode{}

	for _, dev := range devices {
		nodes = append(nodes, dev.makeNode())
	}

	// Loop through each node
	for parentIdx := range nodes {
		parentDepth := len(nodes[parentIdx].Path)

		// Find this node's children
		for childIdx := range nodes {
			if isChild(nodes[parentIdx], nodes[childIdx]) {
				nodes[parentIdx].Children = append(nodes[parentIdx].Children, nodes[childIdx])
			}
		}

		if parentDepth == 0 {
			roots = append(roots, nodes[parentIdx])
		}
	}
	return roots
}

// isChild takes two TreeNodes and returns true if the second is an immediate child of the first.
func isChild(parent TreeNode, child TreeNode) bool {
	// Only looking for immediate children
	if len(parent.Path) != (len(child.Path) - 1) {
		return false
	}

	for i := range parent.Path {
		if parent.Path[i] != child.Path[i] {
			return false
		}
	}

	return true
}

func (d *Device) makeNode() TreeNode {
	return TreeNode{
		Device:   *d,
		Children: []TreeNode{},
	}
}

// Returns sorted slice of Device
// Sorts devices by Bus length, then Path.
func sortDevices(devices []Device) []Device {
	sort.Slice(devices, func(i, j int) bool {
		if devices[i].Bus != devices[j].Bus {
			return devices[i].Bus < devices[j].Bus
		}
		for pathIdx := 0; pathIdx < len(devices[i].Path) && pathIdx < len(devices[j].Path); pathIdx++ {
			if devices[i].Path[pathIdx] != devices[j].Path[pathIdx] {
				return devices[i].Path[pathIdx] < devices[j].Path[pathIdx]
			}
		}
		return len(devices[i].Path) < len(devices[j].Path)
	})
	return devices
}

func addDeviceLog(device Device, logtime time.Time) {
	logs = append(logs, Log{Time: logtime, Text: device.Name, State: device.State})
}

// GetLog returns all stored device logs.
func GetLog() []Log {
	return logs
}
