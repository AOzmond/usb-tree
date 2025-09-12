package lib

import (
	"testing"
	"time"

	"github.com/google/gousb"
	"github.com/stretchr/testify/assert"
)

var (
	device1 = Device{Path: []int{1}, Name: "Device 1", VendorID: "0001", ProductID: "0010", Speed: "High", Bus: 1, State: StateNormal}
	device2 = Device{Path: []int{2}, Name: "Device 2", VendorID: "0002", ProductID: "0020", Speed: "High", Bus: 1, State: StateNormal}
	device3 = Device{Path: []int{3}, Name: "Device 3", VendorID: "0003", ProductID: "0030", Speed: "High", Bus: 1, State: StateNormal}
	device4 = Device{Path: []int{}, Name: "Root", VendorID: "0004", ProductID: "0040", Speed: "High", Bus: 1, State: StateNormal}
	device5 = Device{Path: []int{1}, Name: "Child", VendorID: "0005", ProductID: "0050", Speed: "High", Bus: 1, State: StateNormal}
	device6 = Device{Path: []int{1, 2}, Name: "Grandchild", VendorID: "0006", ProductID: "0060", Speed: "High", Bus: 1, State: StateNormal}
)

var allDevices = []Device{device1, device2, device3, device4, device5}

// mockDesc returns a mock gousb.DeviceDesc for testing descToDevice.
func mockDesc() gousb.DeviceDesc {
	return gousb.DeviceDesc{
		Bus:     1,
		Path:    []int{2, 3},
		Vendor:  gousb.ID(uint16(0x1d6b)),
		Product: gousb.ID(uint16(0x0003)),
		Speed:   gousb.SpeedHigh,
	}
}

// fakeRefresh resets the cached devices and last diff state to the input devices.
func fakeRefresh(newDevices []Device) {
	cachedDevices = newDevices
	deviceDiff(cachedDevices, time.Now())
}

// hasState returns true if a specific state is found within a list of Devices.
func hasState(devs []Device, state LogState) bool {
	for _, d := range devs {
		if d.State == state {
			return true
		}
	}

	return false
}

func TestDescToDevice(t *testing.T) {
	desc := mockDesc()
	dev := descToDevice(desc)
	assert.Equal(t, "3.0 root hub (Linux Foundation)", dev.Name)
	assert.Equal(t, "1d6b", dev.VendorID)
	assert.Equal(t, "0003", dev.ProductID)
	assert.Equal(t, "high", dev.Speed)
	assert.Equal(t, 1, dev.Bus)
}

func TestDeviceDiff_Add(t *testing.T) {
	fakeRefresh([]Device{device1, device2})
	changed, merged := deviceDiff([]Device{device1, device2, device3}, time.Now())
	assert.Len(t, merged, 3)
	assert.True(t, changed)
	assert.True(t, hasState(merged, StateAdded))

	found := false
	for _, dev := range merged {
		if dev.Name == device3.Name {
			found = true
			assert.Equal(t, StateAdded, dev.State)
		}
	}

	assert.True(t, found, "could not find device3 in merged devices")
}

func TestDeviceDiff_Remove(t *testing.T) {
	fakeRefresh([]Device{device1, device2, device3})
	changed, merged := deviceDiff([]Device{device1, device2}, time.Now())
	assert.Len(t, merged, 3)
	assert.True(t, changed)
	assert.True(t, hasState(merged, StateRemoved))

	found := false
	for _, dev := range merged {
		if dev.Name == device3.Name {
			found = true
			assert.Equal(t, StateRemoved, dev.State)
		}
	}

	assert.True(t, found, "could not find device3 in merged devices")
}

func TestDeviceDiff_NoChange(t *testing.T) {
	fakeRefresh([]Device{device1})
	changed, merged := deviceDiff([]Device{device1}, time.Now())
	assert.Len(t, merged, 1)
	assert.False(t, changed)
}

func TestDeviceDiff_AddAndRemove(t *testing.T) {
	fakeRefresh([]Device{device1, device2})
	changed, merged := deviceDiff([]Device{device2, device3}, time.Now())
	assert.True(t, changed)
	assert.True(t, hasState(merged, StateAdded))
	assert.True(t, hasState(merged, StateRemoved))

	foundRemoved := false
	for _, dev := range merged {
		if dev.Name == device1.Name {
			foundRemoved = true
			assert.Equal(t, StateRemoved, dev.State)
		}
	}

	assert.True(t, foundRemoved, "could not find removed device1")

	foundAdded := false
	for _, dev := range merged {
		if dev.Name == device3.Name {
			foundAdded = true
			assert.Equal(t, StateAdded, dev.State)
		}
	}

	assert.True(t, foundAdded, "could not find added device3")
}

func TestDeviceDiff_AddThenRemove(t *testing.T) {
	fakeRefresh([]Device{device1})
	deviceDiff([]Device{device1, device2}, time.Now())
	changed, merged := deviceDiff([]Device{device1}, time.Now())
	assert.True(t, changed)
	assert.False(t, hasState(merged, StateAdded))
	assert.False(t, hasState(merged, StateRemoved))

	found := false
	for _, dev := range merged {
		if dev.Name == device2.Name {
			found = true
		}
	}

	assert.False(t, found, "device2 should not be found after removal")
}

func TestIsChild(t *testing.T) {
	parent := TreeNode{
		Device: Device{
			Path: []int{1, 2},
			Bus:  1,
		},
	}
	child := TreeNode{
		Device: Device{
			Path: []int{1, 2, 3},
			Bus:  1,
		},
	}
	notChild := TreeNode{
		Device: Device{
			Path: []int{1, 3, 4},
			Bus:  1,
		},
	}
	differentBus := TreeNode{
		Device: Device{
			Path: []int{1, 2, 3},
			Bus:  2,
		},
	}

	assert.True(t, isChild(parent, child))
	assert.False(t, isChild(parent, notChild))
	assert.False(t, isChild(parent, parent))
	assert.False(t, isChild(child, parent))
	assert.False(t, isChild(child, notChild))
	assert.False(t, isChild(parent, differentBus))
}

func TestBuildDeviceTree(t *testing.T) {
	tree := BuildDeviceTree([]Device{device4, device5, device6})
	assert.Equal(t, 1, len(tree), "expected 1 root node")
	assert.Equal(t, 1, len(tree[0].Children), "expected 1 child of root")
	assert.Equal(t, device4, tree[0].Device, "root name mismatch")
	assert.Equal(t, device5, tree[0].Children[0].Device, "child name mismatch")
	assert.Equal(t, device6, tree[0].Children[0].Children[0].Device, "grandchild name mismatch")
}

func TestSortDeviceSlice(t *testing.T) {
	sorted := sortDevices(allDevices)
	want := []Device{device4, device1, device5, device2, device3}

	assert.Equal(t, len(want), len(sorted), "sorted slice length mismatch")
	for i := range want {
		assert.Equal(t, want[i].Path, sorted[i].Path, "device path mismatch at index %d", i)
		assert.Equal(t, want[i].Name, sorted[i].Name, "device name mismatch at index %d", i)
	}
}

func TestAddDeviceLogAndGetLog(t *testing.T) {
	logs = nil
	d := Device{Name: "TestLog", State: StateAdded}
	logtime := time.Now()
	addDeviceLog(d, logtime)
	got := GetLog()
	assert.NotEmpty(t, got, "log should not be empty")
	assert.Equal(t, "TestLog", got[0].Text)
	assert.Equal(t, StateAdded, got[0].State)
}

func TestDeviceDiffProducesLog(t *testing.T) {
	fakeRefresh([]Device{device1, device2})
	logtime := time.Now()
	logs = nil
	deviceDiff([]Device{device1, device2, device3}, logtime)
	logsAfter := GetLog()
	assert.Len(t, logsAfter, 1, "expected exactly one new log entry")
	assert.Equal(t, device3.Name, logsAfter[0].Text, "log entry device name should match device3")
	assert.Equal(t, StateAdded, logsAfter[0].State, "log entry state should be StateAdded")
}
