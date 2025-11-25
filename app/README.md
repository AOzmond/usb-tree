# USB Tree - GUI Application

The USB Tree GUI provides an intuitive interface for viewing connected USB devices, monitoring hot-plug
events, and accessing detailed device information.

Built with [Wails v2](https://wails.io/), it combines a native Go backend with a modern Svelte frontend for
optimal performance and user experience.

See the Website at [usb-tree.github.io](https://usb-tree.github.io)

---

## Tech Stack

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Svelte](https://img.shields.io/badge/Svelte-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)](https://svelte.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Wails](https://img.shields.io/badge/Wails-FF6C37?style=for-the-badge&logo=wails&logoColor=white)](https://wails.io/)
[![Bun](https://img.shields.io/badge/Bun-000000?style=for-the-badge&logo=bun&logoColor=white)](https://bun.sh/)
[![Vite](https://img.shields.io/badge/Vite-646CFF?style=for-the-badge&logo=vite&logoColor=white)](https://vitejs.dev/)
[![SCSS](https://img.shields.io/badge/SCSS-CC6699?style=for-the-badge&logo=sass&logoColor=white)](https://sass-lang.com/)
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

- **Real-Time USB Device Tree**: Visualize connected USB devices in a hierarchical tree structure
  - Shows manufacturer and product names
  - Shows the speed the device is using
- **Hot-Plug Detection**: Automatically detect when devices are connected or disconnected
  - Device tree with visual indicators for changes
  - Log of device disconnects and connects
- **Device Details**: Device tooltip shows comprehensive information including:
  - Vendor ID and Product ID
  - Device bus information
  - Click to search in an online device database
- **Cross-Platform**: Native support for Linux (x86-64 and ARM64) and Windows (x86-64)
- **Modern UI**: Clean, responsive interface built with Svelte and Carbon Design System with dark mode support

## Demo

<img src="/images/output.gif" alt="Screen capture showing the functions" width="500">

## Installation

### Linux

#### From Release

1. Download the latest
   [`usb-tree-linux-amd64.tar.gz`](https://github.com/AOzmond/usb-tree/releases/latest/download/usb-tree-linux-amd64.tar.gz)
   or
   [`usb-tree-linux-arm64.tar.gz`](https://github.com/AOzmond/usb-tree/releases/latest/download/usb-tree-linux-arm64.tar.gz)
   or see the [releases](https://github.com/AOzmond/usb-tree/releases) page for all versions.

2. Extract the archive to your desired location:

   ```bash
   tar -xzf usb-tree-linux-*.tar.gz
   ```

3. Run binary:

   ```bash
   ./usb-tree
   ```

#### Arch Linux (AUR)

AUR helper

```bash
# Binary package
yay -S usb-tree-app-bin # or paru etc.
```

or

```bash
# Build from source
yay -S usb-tree-app
```

Manually

```bash
# Binary package
git clone https://aur.archlinux.org/usb-tree-app-bin.git
cd usb-tree-app-bin
makepkg -si
```

or

```bash
# Build from source
git clone https://aur.archlinux.org/usb-tree-app.git
cd usb-tree-app
makepkg -si
```

### Windows

#### Winget

```bash
winget install USBTree.USBTree
```

#### Portable

1. Download the latest
   [`usb-tree-windows-amd64.zip`](https://github.com/AOzmond/usb-tree/releases/latest/download/usb-tree-windows-amd64.zip)
   or see the [releases](https://github.com/AOzmond/usb-tree/releases) page for all versions.
2. Extract the archive to your desired location.
3. Run `usb-tree.exe`.

Note: Ensure `libusb-1.0.dll` is in the same directory as the executable.

#### Installer

Download and run the latest
[`usb-tree-amd64-installer.exe`](https://github.com/AOzmond/usb-tree/releases/latest/download/usb-tree-amd64-installer.exe)
or see the [releases](https://github.com/AOzmond/usb-tree/releases) page for all versions.

### Building from Source

#### Required build tools

- [**Go**: 1.25](https://go.dev/dl/)
- [**Bun**: 1.3.1](https://bun.sh/)
- [**Wails (v2)**: 2.10](https://wails.io/docs/next/gettingstarted/installation)

#### System dependencies

- [**libusb-1.0**](https://libusb.info/)

### Build Steps

Clone the repository and build the application.

```bash
git clone https://github.com/AOzmond/usb-tree.git
cd usb-tree/app
wails build
```

The compiled executable will be in `build/bin`, named `usb-tree` for Linux, and `usb-tree.exe` for Windows.

## License

This project is licensed under the GPL-2.0 License. See the [`LICENSE`](../LICENSE) file for details.

## Author

Alastair Ozmond

[![Looking for Work](https://img.shields.io/badge/hiring-I'm%20looking%20for%20work-blue?style=flat-square)](https://www.linkedin.com/in/alastair-ozmond-108512179)

[LinkedIn](https://www.linkedin.com/in/alastair-ozmond-108512179)

Software Engineer | Full-Stack & Systems Developer
