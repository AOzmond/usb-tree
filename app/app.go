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

// runs on startup to
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// InitFrontend initializes the usb tree library with the app.updateCallback
func (a *App) InitFrontend() {
	lib.Init(a.updateCallback)
}

// Refresh relays refresh request to library, sets updated device tree on frontend
func (a *App) Refresh() []lib.Device {
	return lib.Refresh()
}

func (a *App) Exit() {
	lib.Stop()
}

// TODO better name maybe?
func (a *App) updateCallback(newDevices []lib.Device) {
	tree := lib.BuildDeviceTree(newDevices)
	logs := lib.GetLog()
	runtime.EventsEmit(a.ctx, "treeUpdated", tree)
	runtime.EventsEmit(a.ctx, "logsUpdated", logs)
}
