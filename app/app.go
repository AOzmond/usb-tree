package main

import (
	"context"

	"github.com/AOzmond/usb-tree/lib"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// InitFrontend initializes the usb tree library with the app.updateCallback
func (a *App) InitFrontend() {
	lib.Init(a.updateCallback)
}

// Refresh relays refresh request to library, sets updated device tree on frontend
func (a *App) Refresh() {
	lib.Refresh()
}

// Exit will stop polling of new devices.
func (a *App) Exit(ctx context.Context) {
	lib.Stop()
}

// updateCallback will emit update events on device changes.
func (a *App) updateCallback(newDevices []lib.Device) {
	tree := lib.BuildDeviceTree(newDevices)
	logs := lib.GetLog()
	runtime.EventsEmit(a.ctx, "treeUpdated", tree)
	runtime.EventsEmit(a.ctx, "logsUpdated", logs)
}
