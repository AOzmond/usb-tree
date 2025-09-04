package usb_tree_lib

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
)

const (
	StateNormal  LogState = "normal"
	StateAdded   LogState = "added"
	StateRemoved LogState = "removed"
)

type LogState string

type Device struct {
	Path      []int
	Name      string
	VendorId  string
	ProductId string
	Speed     string
	Bus       int
	State     LogState
}

type TreeNode struct {
	Device
	Children []TreeNode
}

type Log struct {
	Time  time.Time
	Text  string
	State LogState
}

var (
	cachedDevices []Device
	lastMergedMap map[string]Device
	logs          []Log
)

var pollingStop chan bool

func Stop() {
	select {
	case pollingStop <- true:
	default:
	}
}

func Init(onUpdateCallback func([]Device)) []Device {
	pollingStop = make(chan bool)
	Refresh()
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				newDevices := getDevices()
				changed, mergedDevices := deviceDiff(cachedDevices, newDevices)
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

// Resets tree to current state
func Refresh() []Device {
	cachedDevices = getDevices()
	return cachedDevices
}

// returns lists of devices in depth first search order.
func getDevices() []Device {
	ctx := gousb.NewContext()
	defer ctx.Close()
	devices := []Device{}

	_, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		devices = append([]Device{descToDevice(*desc)}, devices...)
		return false
	})
	if err != nil {
		return nil
	}
	return devices
}

// Returns device based on a given DeviceDesc
func descToDevice(desc gousb.DeviceDesc) Device {
	return Device{
		Path:      desc.Path,
		Name:      usbid.Describe(desc),
		VendorId:  fmt.Sprintf("%04x", desc.Vendor),
		ProductId: fmt.Sprintf("%04x", desc.Product),
		Speed:     fmt.Sprintf("%v", desc.Speed),
		Bus:       desc.Bus,
		State:     StateNormal,
	}
}

func deviceDiff(cachedDevices []Device, newDevices []Device) (changed bool, merged []Device) {
	mergedMap := make(map[string]Device)
	changed = false

	deviceKey := func(d Device) string {
		return fmt.Sprintf("%d:%v:%s:%s:%s", d.Bus, d.Path, d.VendorId, d.ProductId, d.Speed)
	}

	// Mark all cachedDevices as removed initially
	for _, device := range cachedDevices {
		device.State = StateRemoved
		mergedMap[deviceKey(device)] = device
	}
	// Reset persisting devices to normal, and add new devices
	for _, newDevice := range newDevices {
		key := deviceKey(newDevice)
		if existingDevice, exists := mergedMap[key]; exists {
			// Device exists, reset its status to normal
			existingDevice.State = StateNormal
			mergedMap[key] = existingDevice
		} else {
			// Device is new, add to mergedMap
			newDevice.State = StateAdded
			mergedMap[key] = newDevice
			changed = true
		}
	}
	// Convert map back into slice and log changes since lastMergedMap
	merged = make([]Device, 0, len(mergedMap))
	for key, device := range mergedMap {
		merged = append(merged, device)
		if _, exists := lastMergedMap[key]; !exists {
			addDeviceLog(device)
		}
	}

	merged = sortDeviceSlice(merged)

	lastMergedMap = mergedMap

	return changed, merged
}

func BuildDeviceTree(devices []Device) []TreeNode {
	roots := []TreeNode{}
	nodes := []TreeNode{}
	for _, dev := range devices {
		nodes = append(nodes, makeNode(dev))
	}
	// Loop through each node
	for _, parent := range nodes {
		parentLen := len(parent.Path)
		if parentLen == 0 {
			roots = append(roots, parent)
		}
		// Find this node's children
		for _, child := range nodes {
			childLen := len(child.Path)
			if parentLen >= childLen {
				break
			}
			if isChild(parent, child) {
				parent.Children = append(parent.Children, child)
			}
		}
	}
	return roots
}

/*
	 Checks parent and child relationship
		Criteria: 	Same Bus
					Path of Child matches path of parent with 1 extra port.
*/
func isChild(parent TreeNode, child TreeNode) bool {
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

func makeNode(device Device) TreeNode {
	return TreeNode{
		Device:   device,
		Children: []TreeNode{},
	}
}

// Returns sorted slice of Device
// Sorts devices by Bus length, then Path.
func sortDeviceSlice(devices []Device) []Device {
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

func addDeviceLog(device Device) {
	newLog := Log{Time: time.Now(), Text: device.Name, State: device.State}
	logs = append(logs, newLog)
}

func GetLog() []Log {
	return logs
}
