package lib

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/gousb"
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
	if dev.Name != "3.0 root hub (Linux Foundation)" {
		t.Errorf(" Name = %s, want 3.0 root hub (Linux Foundation)", dev.Name)
	}
	if dev.VendorID != "1d6b" {
		t.Errorf("VendorID = %s, want 1d6b", dev.VendorID)
	}
	if dev.ProductID != "0003" {
		t.Errorf("ProductID = %s, want 0003", dev.ProductID)
	}
	if dev.Speed != "high" {
		t.Errorf("Speed = %s, want high", dev.Speed)
	}
	if dev.Bus != 1 {
		t.Errorf("Bus = %d, want 1", dev.Bus)
	}
}

func TestDeviceDiff_Add(t *testing.T) {
	fakeRefresh([]Device{device1, device2})
	changed, merged := deviceDiff([]Device{device1, device2, device3}, time.Now())
	if len(merged) != 3 {
		t.Errorf("length of merged = %d, want 3", len(merged))
	}
	if !changed {
		t.Errorf("changed = false, want true")
	}
	if !hasState(merged, StateAdded) {
		t.Errorf("hasState(merged, stateAdded) = false, want true")
	}
	found := false
	for _, dev := range merged {
		if dev.Name == device3.Name {
			found = true
			if dev.State != StateAdded {
				t.Errorf("State = %v, want  %v", dev.State, StateAdded)
			}
		}
	}
	if !found {
		t.Errorf("could not find device3 in merged devices")
	}
}

func TestDeviceDiff_Remove(t *testing.T) {
	fakeRefresh([]Device{device1, device2, device3})
	changed, merged := deviceDiff([]Device{device1, device2}, time.Now())
	if len(merged) != 3 {
		t.Errorf("merged length = %d, want 3", len(merged))
	}
	if !changed {
		t.Errorf("changed = false, want true")
	}
	if !hasState(merged, StateRemoved) {
		t.Errorf("hasState(merged, StateRemoved) = false, want true")
	}

	found := false
	for _, dev := range merged {
		if dev.Name == device3.Name {
			found = true
			if dev.State != StateRemoved {
				t.Errorf("State = %v, want %v", dev.State, StateRemoved)
			}
		}
	}
	if !found {
		t.Errorf("could not find device3 in merged devices")
	}
}

func TestDeviceDiff_NoChange(t *testing.T) {
	fakeRefresh([]Device{device1})
	changed, merged := deviceDiff([]Device{device1}, time.Now())
	if len(merged) != 1 {
		t.Errorf("length of merged = %d, want 1", len(merged))
	}
	if changed {
		t.Errorf("changed = true, want false")
	}
}

func TestDeviceDiff_AddAndRemove(t *testing.T) {
	fakeRefresh([]Device{device1, device2})
	changed, merged := deviceDiff([]Device{device2, device3}, time.Now())
	if !changed {
		t.Errorf("changed = false, want true")
	}
	if !hasState(merged, StateAdded) {
		t.Errorf("hasState(merged, StateAdded) = false, want true")
	}
	if !hasState(merged, StateRemoved) {
		t.Errorf("hasState(merged, StateRemoved) = false, want true")
	}

	found := false
	for _, dev := range merged {
		if dev.Name == device1.Name {
			found = true
			if dev.State != StateRemoved {
				t.Errorf("State = %v, want %v", dev.State, StateRemoved)
			}
		}
	}
	if !found {
		t.Errorf("found = false, want true.")
	}

	found = false
	for _, dev := range merged {
		if dev.Name == device3.Name {
			found = true
			if dev.State != StateAdded {
				t.Errorf("State = %v, want %v", dev.State, StateAdded)
			}
		}
	}
	if !found {
		t.Errorf("found = false, want true.")
	}
}

func TestDeviceDiff_AddThenRemove(t *testing.T) {
	fakeRefresh([]Device{device1})
	deviceDiff([]Device{device1, device2}, time.Now())
	changed, merged := deviceDiff([]Device{device1}, time.Now())
	if !changed {
		t.Errorf("changed = false, want true")
	}
	if hasState(merged, StateAdded) {
		t.Errorf("hasState(merged, StateAdded) = true, want false")
	}
	if hasState(merged, StateRemoved) {
		t.Errorf("hasState(merged, StateRemoved) = true, want false")
	}

	found := false
	for _, dev := range merged {
		if dev.Name == device2.Name {
			found = true
		}
	}
	if found {
		t.Errorf("found = true, want false.")
	}
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

	if !isChild(parent, child) {
		t.Errorf("isChild(parent, child) = false, want true")
	}

	if isChild(parent, notChild) {
		t.Errorf("isChild(parent, notChild) = true, want false")
	}

	if isChild(parent, parent) {
		t.Errorf("isChild(parent, parent) = true, want false")
	}

	if isChild(child, parent) {
		t.Errorf("isChild(child, parent) = true, want false")
	}

	if isChild(child, notChild) {
		t.Errorf("isChild(child, notChild) = true, want false")
	}
	if isChild(parent, differentBus) {
		t.Errorf("isChild(child, notChild) = true, want false")
	}
}

func TestBuildDeviceTree(t *testing.T) {
	tree := BuildDeviceTree([]Device{device4, device5, device6})
	if len(tree) != 1 {
		t.Errorf("number of roots = %d, wanted 1", len(tree))
	}
	if len(tree) > 0 && len(tree[0].Children) != 1 {
		t.Errorf("length of roots children = %d, expected 1", len(tree[0].Children))
	}
	if tree[0].Name != device4.Name {
		t.Errorf("Root wrong")
	}
	if tree[0].Children[0].Name != device5.Name {
		t.Errorf("child wrong")
	}
	if tree[0].Children[0].Children[0].Name != device6.Name {
		t.Errorf("Wrong Grandchild")
	}
}

func TestSortDeviceSlice(t *testing.T) {
	sorted := sortDevices(allDevices)
	want := []Device{device4, device1, device5, device2, device3}

	for i := range want {
		if !reflect.DeepEqual(sorted[i].Path, want[i].Path) {
			t.Errorf("got %+v want %+v", sorted[i].Name, want[i].Name)
		}
	}
}

func TestAddDeviceLogAndGetLog(t *testing.T) {
	logs = nil
	d := Device{Name: "TestLog", State: StateAdded}
	logtime := time.Now()
	addDeviceLog(d, logtime)
	got := GetLog()
	if len(got) == 0 || got[0].Text != "TestLog" || got[0].State != StateAdded {
		t.Errorf("got %+v", got)
	}
}

func TestDeviceDiffProducesLog(t *testing.T) {
	logs = nil
	fakeRefresh([]Device{device1, device2})
	logtime := time.Now()
	logsBefore := GetLog()
	deviceDiff([]Device{device1, device2, device3}, logtime)
	logsAfter := GetLog()
	if (len(logsAfter) - len(logsBefore)) != 1 {
		t.Fatalf("logs empty, want 3")
	}
	found := false
	for _, l := range logsAfter {
		if l.Text == device3.Name && l.State == StateAdded && l.Time.Equal(logtime) {
			found = true
		}
	}
	if !found {
		t.Errorf("expected log for added device3, got %+v", logsAfter)
	}
}
