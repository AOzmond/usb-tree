# USB Tree - GUI Application

The USB Tree GUI provides an intuitive interface for viewing connected USB devices, monitoring hot-plug
events, and accessing detailed device information. Built with [Wails v2](https://wails.io/), it combines a
native Go backend with a modern Svelte frontend for optimal performance and user experience.

[![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)](https://kernel.org/)
[![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)](https://microsoft.com/windows)

---

## Table of Contents

- [Features](#features)
- [Demo](#demo)
- [Installation](#installation)
  - [Linux](#linux)
  - [Windows](#windows)
  - [Building from Source](#building-from-source)
  - [Build Steps](#build-steps)
- [License](#license)
- [Author](#author)

## Features

- **Real-Time USB Device Tree**: Visualize all connected USB devices in a hierarchical tree structure
- **Hot-Plug Detection**: Automatically detect when devices are connected or disconnected
  - changes are reflected in the tree structure with color coding and icons
  - all changes are logged to a log section at the bottom of the app
  - Device tree can be reset to the current state by clicking the "Refresh" button
- **Device Details**: View comprehensive information including:
  - Vendor ID and Product ID
  - Manufacturer and product names
  - Device bus information
  - USB version and device speed
  - Clicking a device in the tree will open a database search for more information.
- **Event Logging**: Track all device connection and disconnection events with timestamps in the lower portion
  of the app.
- **Cross-Platform**: Native support for Linux and Windows
- **Modern UI**: Clean, responsive interface built with Svelte and Carbon Design System
  - Dark mode supported

## Demo:

<img src="/images/output.gif" alt="Demo" width="500" >

## Installation

### Linux

#### From Release

1. Download the latest `usb-tree-linux-amd64.tar.gz` from the
   [releases page](https://github.com/AOzmond/usb-tree/releases)
2. Extract the archive to your desired location:
   ```bash
   tar -xzf usb-tree-linux-amd64.tar.gz
   ```
3. Run binary:
   ```bash
   ./usb-tree
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

### Building from Source

#### Required build tools:

- [**Go**: 1.25 ](https://go.dev/dl/)
- [**Bun**: 1.3.1](https://bun.sh/)
- [**Wails (v2)**: 2.10](https://wails.io/docs/next/gettingstarted/installation)
  - see wails install page for OS specific instruction.

#### System dependencies:

- [**libusb-1.0**](https://libusb.info/)

### Build Steps

1. **Clone the repository:**

   ```bash
   git clone https://github.com/AOzmond/usb-tree.git
   cd usb-tree/app
   ```

2. **Install Wails CLI:**

   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

3. **Build the application:**

   ```bash
   wails build -clean
   ```

4. **Locate the binary:**

   The compiled executable will be in:

- **Linux**: `build/bin/usb-tree`
- **Windows**: `build/bin/usb-tree.exe`

## License

This project is licensed under the GPL-2.0 License. See the LICENSE file for details.

## Author

Alastair Ozmond

[![Looking for Work](https://img.shields.io/badge/hiring-I'm%20looking%20for%20work-blue?style=flat-square)](https://www.linkedin.com/in/alastair-ozmond-108512179)

[LinkedIn](https://www.linkedin.com/in/alastair-ozmond-108512179)

Software Engineer | Full-Stack & Systems Developer
