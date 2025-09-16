package lib

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
	Path      []int    `json:"path"`
	Name      string   `json:"name"`
	VendorID  string   `json:"vendorId"`
	ProductID string   `json:"productId"`
	Speed     string   `json:"speed"`
	Bus       int      `json:"bus"`
	State     LogState `json:"state"`
}

// TreeNode represents a Device and its children for building tree structures.
type TreeNode struct {
	Device   `json:"device"`
	Children []*TreeNode `json:"children"`
}

// Log represents a change in a Device.
type Log struct {
	Time  time.Time
	Text  string
	Speed string
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
	onUpdateCallback(Refresh())
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

// returns lists of devices.
func getDevices() (time.Time, []Device) {
	ctx := gousb.NewContext()
	defer func() {
		if err := ctx.Close(); err != nil {
			log.Printf("error closing USB context: %v\n", err)
		}
	}()

	devices := []Device{}

	_, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		devices = append(devices, descToDevice(*desc))
		return false
	})
	if err != nil {
		log.Printf("Issue accessing USB Devices: %v", err)
		return time.Now(), []Device{}
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

func (d *Device) key() string {
	return fmt.Sprintf("%d:%v:%s:%s:%s", d.Bus, d.Path, d.VendorID, d.ProductID, d.Speed)
}

func deviceDiff(newDevices []Device, logtime time.Time) (changed bool, merged []Device) {
	mergedMap := make(map[string]Device)
	changed = false

	// Mark all cachedDevices as removed initially
	for _, device := range cachedDevices {
		device.State = StateRemoved
		mergedMap[device.key()] = device
	}

	// Reset persisting devices to normal, and add new devices
	for _, newDevice := range newDevices {
		key := newDevice.key()
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

	// Search for removed devices to update changed.
	for key := range lastMergedMap {
		if _, exists := mergedMap[key]; !exists {
			changed = true
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

// BuildDeviceTree converts a device list to a device tree
func BuildDeviceTree(devices []Device) []*TreeNode {
	roots := []*TreeNode{}
	nodes := []*TreeNode{}

	for _, dev := range devices {
		newNode := dev.treeNode()
		nodes = append(nodes, &(newNode))
	}

	// Loop through each node to assign children
	for parentIdx := len(nodes) - 1; parentIdx >= 0; parentIdx-- {
		// Find this node's children
		for childIdx := range nodes {

			parent := nodes[parentIdx]
			child := nodes[childIdx]
			if isChild(*parent, *child) {
				parent.Children = append(parent.Children, child)
			}
		}
	}

	for _, node := range nodes {
		if len(node.Path) == 0 {
			roots = append(roots, node)
		}
	}

	return roots
}

// isChild checks if maybeChild is an immediate child of the parent
func isChild(parent TreeNode, maybeChild TreeNode) bool {
	if parent.Bus != maybeChild.Bus {
		return false
	}
	// Only looking for immediate children
	if len(parent.Path) != (len(maybeChild.Path) - 1) {
		return false
	}

	for i := range parent.Path {
		if parent.Path[i] != maybeChild.Path[i] {
			return false
		}
	}

	return true
}

func (d *Device) treeNode() TreeNode {
	return TreeNode{
		Device:   *d,
		Children: []*TreeNode{},
	}
}

// sortDevices sorts devices consistently by their path
func sortDevices(devices []Device) []Device {
	sort.Slice(devices, func(i, j int) bool {
		if devices[i].Bus != devices[j].Bus {
			return devices[i].Bus < devices[j].Bus
		}
		flatten := func(path []int) string {
			s := ""
			for _, p := range path {
				s += fmt.Sprintf("%04d-", p)
			}
			return s
		}

		return flatten(devices[i].Path) < flatten(devices[j].Path)
	})

	return devices
}

func addDeviceLog(device Device, logtime time.Time) {
	logs = append(logs, Log{Time: logtime, Text: device.Name, State: device.State, Speed: device.Speed})
}

// GetLog returns all stored device logs.
func GetLog() []Log {
	return logs
}
