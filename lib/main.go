package toollib

import (
	"fmt"
	"time"

	"github.com/google/gousb"
)

type Device struct {
	Path      []int
	Name      string
	VendorId  string
	ProductId string
	Speed     string
	Bus       int
	State     LogState
}

type LogState string

const (
	StateNormal  LogState = "normal"
	StateAdded   LogState = "added"
	StateRemoved LogState = "removed"
)

type TreeNode struct {
	Device
	Children []TreeNode
}

type Log struct {
	Time  time.Time
	Text  string
	State LogState
}

var cachedDevices []Device
var logs []Log

var pollingStop chan bool

func Stop() {
	select {
	case pollingStop <- true:
	default:
	}
}

func Init(onUpdateCallback func([]Device)) []Device {
	pollingStop = make(chan bool)
	cachedDevices = Refresh()
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				newDevices := getDevices()
				changed, cachedDevices := deviceDiff(cachedDevices, newDevices)
				if changed {
					onUpdateCallback(cachedDevices)
				}
			case <-pollingStop:
				close(pollingStop)
				return

			}
		}
	}()

	return cachedDevices
}

// resets cachedDevices to current config and returns Devices
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

// returns device based on a given DeviceDesc
func descToDevice(desc gousb.DeviceDesc) Device {
	return Device{
		Path:      desc.Path,
		Name:      desc.String(),
		VendorId:  fmt.Sprintf("%04x", desc.Vendor),
		ProductId: fmt.Sprintf("%04x", desc.Product),
		Speed:     fmt.Sprintf("%v", desc.Speed),
		Bus:       desc.Bus,
		State:     StateNormal,
	}
}

// TODO
func deviceDiff(cachedDevices []Device, newDevices []Device) (changed bool,merged []Device) {
	//TODO make sure return is sorted Depth first

	return changed, merged
}

func BuildDeviceTree(devices []Device) []TreeNode {
	result := []TreeNode{}
	//TODO sort devices, maybe not (depth first on insert)
	//TODO make the appends references instead of copies (?)
	nodes := []TreeNode{}
	for _, dev := range devices {
		nodes = append(nodes, makeNode(dev))
	}
	//loop through each parent
	for pIndex, parent := range nodes {
		pLength := len(parent.Path)
		if pLength == 0 {
			result = append(result, parent)
		}
		//loop through each possible child
		for cIndex := pIndex + 1; cIndex < len(nodes); cIndex++ {
			cLength := len(parent.Path)
			if pLength >= cLength {
				break
			}
			if isChild(parent, nodes[cIndex]) {
				parent.Children = append(parent.Children, nodes[cIndex])
			}
		}
	}

	return result
}

/*
	 Checks parent and child relationship
		Criteria: 	Same BUS
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

func GetLog() []Log {
	return logs
}

// TODO remove this probably
func HelloWorld() string {
	return "Hello World!"
}
