package usb_tree_lib

import (
	"reflect"
	"testing"

	"github.com/google/gousb"
)

func TestDescToDevice(t *testing.T) {
	desc := mockDesc()
	dev := descToDevice(desc)
	if dev.Name != "3.0 root hub (Linux Foundation)" {
		t.Errorf("descToDevice failed: Name = %s, 3.0 root hub (Linux Foundation)", dev.Name)
	}
	if dev.VendorId != "1d6b" {
		t.Errorf("descToDevice failed: VendorId = %s, want 1d6b", dev.VendorId)
	}
	if dev.ProductId != "0003" {
		t.Errorf("descToDevice failed: ProductId = %s, want 0003", dev.ProductId)
	}
	if dev.Speed != "high" {
		t.Errorf("descToDevice failed: Speed = %s, want high", dev.Speed)
	}
	if dev.Bus != 1 {
		t.Errorf("descToDevice failed: Bus = %d, want 1", dev.Bus)
	}
}

func TestDeviceDiff_Add(t *testing.T) {
	d1 := Device{Path: []int{1}, Name: "A", VendorId: "v1", ProductId: "p1", Speed: "High", Bus: 1, State: StateNormal}
	d2 := Device{Path: []int{2}, Name: "B", VendorId: "v2", ProductId: "p2", Speed: "High", Bus: 1, State: StateNormal}
	d3 := Device{Path: []int{3}, Name: "C", VendorId: "v3", ProductId: "p3", Speed: "High", Bus: 1, State: StateNormal}
	mockRefresh([]Device{d1, d2})
	changed, merged := deviceDiff([]Device{d1, d2, d3})
	if len(merged) != 3 {
		t.Errorf("deviceDiff add failed: expected 3 devices, got %d", len(merged))
	}
	if !changed {
		t.Errorf("deviceDiff add failed: expected changed=true, got false")
	}
	if !hasState(merged, StateAdded) {
		t.Errorf("deviceDiff add failed: expected StateAdded in merged devices")
	}
}

func TestDeviceDiff_Remove(t *testing.T) {
	d1 := Device{Path: []int{1}, Name: "A", VendorId: "v1", ProductId: "p1", Speed: "High", Bus: 1, State: StateNormal}
	d2 := Device{Path: []int{2}, Name: "B", VendorId: "v2", ProductId: "p2", Speed: "High", Bus: 1, State: StateNormal}
	d3 := Device{Path: []int{3}, Name: "C", VendorId: "v3", ProductId: "p3", Speed: "High", Bus: 1, State: StateNormal}
	mockRefresh([]Device{d1, d2, d3})
	changed, merged := deviceDiff([]Device{d1})
	if len(merged) != 3 {
		t.Errorf("deviceDiff remove failed: expected 2 devices, got %d", len(merged))
	}
	if !changed {
		t.Errorf("deviceDiff remove failed: expected changed=true, got false")
	}
	if !hasState(merged, StateRemoved) {
		t.Errorf("deviceDiff remove failed: expected StateRemoved in merged devices")
	}
}

func TestDeviceDiff_NoChange(t *testing.T) {
	d1 := Device{Path: []int{1}, Name: "A", VendorId: "v1", ProductId: "p1", Speed: "High", Bus: 1, State: StateNormal}
	mockRefresh([]Device{d1})
	changed, merged := deviceDiff([]Device{d1})
	if len(merged) != 1 {
		t.Errorf("deviceDiff no change failed: expected 1 device, got %d", len(merged))
	}
	if changed {
		t.Errorf("deviceDiff no change failed: expected changed=false, got true")
	}
}

func TestDeviceDiff_AddAndRemove(t *testing.T) {
	d1 := Device{Path: []int{1}, Name: "A", VendorId: "v1", ProductId: "p1", Speed: "High", Bus: 1, State: StateNormal}
	d2 := Device{Path: []int{2}, Name: "B", VendorId: "v2", ProductId: "p2", Speed: "High", Bus: 1, State: StateNormal}
	d3 := Device{Path: []int{3}, Name: "C", VendorId: "v3", ProductId: "p3", Speed: "High", Bus: 1, State: StateNormal}
	mockRefresh([]Device{d1, d2})
	changed, merged := deviceDiff([]Device{d2, d3})
	if !changed {
		t.Errorf("deviceDiff add+remove failed: expected changed=true, got false")
	}
	if !hasState(merged, StateAdded) {
		t.Errorf("deviceDiff add+remove failed: expected StateAdded in merged devices")
	}
	if !hasState(merged, StateRemoved) {
		t.Errorf("deviceDiff add+remove failed: expected StateRemoved in merged devices")
	}
}

func hasState(devs []Device, state LogState) bool {
	for _, d := range devs {
		if d.State == state {
			return true
		}
	}
	return false
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
	// Should be true: child path extends parent path by one element and matches up to parent length
	if !isChild(parent, child) {
		t.Errorf("isChild failed: expected true for child %+v of parent %+v", child, parent)
	}
	// Should be false: notChild path does not match parent path
	if isChild(parent, notChild) {
		t.Errorf("isChild failed: expected false for notChild %+v of parent %+v", notChild, parent)
	}
	// Should be false: same path length
	if isChild(parent, parent) {
		t.Errorf("isChild failed: expected false for same node")
	}
	// Should be false: child path is shorter
	if isChild(child, parent) {
		t.Errorf("isChild failed: expected false when parent is actually child")
	}
}

func TestBuildDeviceTreeAndIsChild(t *testing.T) {
	d1 := Device{Path: []int{}, Name: "Root", VendorId: "v1", ProductId: "p1", Speed: "High", Bus: 1}
	d2 := Device{Path: []int{1}, Name: "Child", VendorId: "v2", ProductId: "p2", Speed: "High", Bus: 1}
	d3 := Device{Path: []int{1, 2}, Name: "Grandchild", VendorId: "v3", ProductId: "p3", Speed: "High", Bus: 1}
	tree := BuildDeviceTree([]Device{d1, d2, d3})
	if len(tree) != 1 {
		t.Errorf("BuildDeviceTree failed: expected 1 root, got %d, tree: %+v", len(tree), tree)
	}
	if len(tree) > 0 && len(tree[0].Children) != 1 {
		t.Errorf("BuildDeviceTree failed: expected 1 child for root, got %d, tree: %+v", len(tree[0].Children), tree)
	}
}

func TestSortDeviceSlice(t *testing.T) {
	d1 := Device{Path: []int{2}, Name: "B", VendorId: "v2", ProductId: "p2", Speed: "High", Bus: 2}
	d2 := Device{Path: []int{1}, Name: "A", VendorId: "v1", ProductId: "p1", Speed: "High", Bus: 1}
	d3 := Device{Path: []int{1, 2}, Name: "C", VendorId: "v3", ProductId: "p3", Speed: "High", Bus: 1}
	d4 := Device{Path: []int{}, Name: "Root", VendorId: "v4", ProductId: "p4", Speed: "High", Bus: 1}
	d5 := Device{Path: []int{2}, Name: "Child", VendorId: "v5", ProductId: "p5", Speed: "High", Bus: 1}
	sorted := sortDeviceSlice([]Device{d1, d2, d3, d4, d5})
	want := []Device{d4, d2, d3, d5, d1}
	for i := range want {
		if !reflect.DeepEqual(sorted[i].Path, want[i].Path) {
			t.Errorf("sortDeviceSlice failed: got %+v want %+v", sorted, want)
		}
	}
}

func TestAddDeviceLogAndGetLog(t *testing.T) {
	logs = nil
	d := Device{Name: "TestLog", State: StateAdded}
	addDeviceLog(d)
	got := GetLog()
	if len(got) == 0 || got[0].Text != "TestLog" || got[0].State != StateAdded {
		t.Errorf("addDeviceLog/GetLog failed: got %+v", got)
	}
}

// mockDesc returns a fake gousb.DeviceDesc for testing descToDevice
func mockDesc() gousb.DeviceDesc {
	return gousb.DeviceDesc{
		Bus:     1,
		Path:    []int{2, 3},
		Vendor:  gousb.ID(uint16(0x1d6b)),
		Product: gousb.ID(uint16(0x0003)),
		Speed:   gousb.SpeedHigh,
	}
}

func mockRefresh(newDevices []Device) {
	cachedDevices = newDevices
	deviceDiff(cachedDevices)
}

// func debug(t *testing.T) {
// 	t.Errorf("CachedDevices : %v\n\n", cachedDevices)
// 	t.Errorf("Last diff map : %v", lastMergedMap)

// }
