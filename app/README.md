# USB Tree - GUI Application

A modern, cross-platform desktop application for visualizing and monitoring USB device topology in real-time.

## Overview

The USB Tree GUI provides an intuitive interface for viewing connected USB devices, monitoring hot-plug
events, and accessing detailed device information. Built with [Wails v2](https://wails.io/), it combines a
native Go backend with a modern Svelte frontend for optimal performance and user experience.

## Features

- **Real-Time USB Device Tree**: Visualize all connected USB devices in a hierarchical tree structure
- **Hot-Plug Detection**: Automatically detect when devices are connected or disconnected
- **Device Details**: View comprehensive information including:
  - Vendor ID and Product ID
  - Manufacturer and product names
  - Device path and bus information
  - USB version and device speed
- **Event Logging**: Track all device connection and disconnection events with timestamps in the lower portion
  of the app.
- **Cross-Platform**: Native support for Linux and Windows
- **Modern UI**: Clean, responsive interface built with Svelte and Carbon Design System

## Installation

### Linux

#### From Release

1. Download the latest `usb-tree-linux-amd64.tar.gz` from the
   [releases page](https://github.com/AOzmond/usb-tree/releases)
2. Extract the archive:
   ```bash
   tar -xzf usb-tree-linux-amd64.tar.gz
   ```
3. Install the files:
   ```bash
   sudo cp usb-tree /usr/local/bin/
   sudo cp usb-tree.desktop /usr/share/applications/
   sudo cp usb-tree.png /usr/share/pixmaps/
   ```

#### Arch Linux (AUR)

```bash
# Binary package
yay -S usb-tree-app-bin

# Or build from source
yay -S usb-tree-app
```

### Windows

1. Download the latest `usb-tree-windows-amd64.zip` from the release page
2. Extract the archive to your desired location
3. Run `usb-tree.exe` Note: Ensure libusb-1.0.dll is in the same directory as the executable.

## License

This project is licensed under the GPL-2.0 License. See the LICENSE file for details.
