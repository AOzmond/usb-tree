# USB Tree

[![Looking for Work](https://img.shields.io/badge/hiring-I'm%20looking%20for%20work-blue?style=flat-square)](https://www.linkedin.com/in/alastair-ozmond-108512179)
[![Build Status](https://github.com/AOzmond/usb-tree/actions/workflows/release.yml/badge.svg)](https://github.com/AOzmond/usb-tree/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/AOzmond/usb-tree/app)](https://goreportcard.com/report/github.com/AOzmond/usb-tree/app)
[![License](https://img.shields.io/github/license/AOzmond/usb-tree)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/AOzmond/usb-tree)](https://github.com/AOzmond/usb-tree/releases)
[![Platform](https://img.shields.io/badge/platform-linux%20|%20windows-lightgrey.svg)](#)
[![Code style: gofumpt](https://img.shields.io/badge/code%20style-gofumpt-blue.svg)](https://github.com/mvdan/gofumpt)

See Website at [usb-tree.github.io](https://usb-tree.github.io)

---

## Overview

**USB Tree** is a cross-platform desktop and command-line tool that visualizes, logs, and monitors USB devices
connected to your computer in real-time. The desktop application is built in **Go** with **Wails v2**,
combining a native Go backend with a modern **Svelte + TypeScript** frontend. The command-line tool is still
under development, but will feature the same Go backend libraries.

Use USB Tree to:

- View all currently connected USB devices in a hierarchical “tree” view
- Log device connect/disconnect events with timestamps
- Display vendor, and product information
- Compare changes across sessions for debugging or audit purposes

The tool is designed for anyone looking for an easy-to-use, lightweight, and portable solution for monitoring
USB devices.

---

## Features

- **Cross-Platform:** Runs on Linux and Windows
- **Dual Interface:** CLI and GUI built from a shared Go library
- **Real-Time Monitoring:** Detect hot-plug events and changes
- **Device Metadata:** Vendor ID, Product ID, Path, and Bus info
- **Lightweight & Fast:** Written in Go for efficiency and portability

---

## Demo

<img src="images/output.gif" alt="Demo" width="500" >

---

## GUI APP

See [GUI documentation](app/README.md) for usage details.

## CLI ( Under Development )

See [CLI documentation](cli/README.md) for usage details.

---

## License

This project is licensed under the GPL-2.0 License.

---

## Author

Alastair Ozmond

[![Looking for Work](https://img.shields.io/badge/hiring-I'm%20looking%20for%20work-blue?style=flat-square)](https://www.linkedin.com/in/alastair-ozmond-108512179)

[LinkedIn](https://www.linkedin.com/in/alastair-ozmond-108512179)

Software Engineer | Full-Stack & Systems Developer
