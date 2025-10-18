# USB Tree
[![Build Status](https://github.com/AOzmond/usb-tree/actions/workflows/release.yml/badge.svg)](https://github.com/AOzmond/usb-tree/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/AOzmond/usb-tree/app)](https://goreportcard.com/report/github.com/AOzmond/usb-tree)
[![License](https://img.shields.io/github/license/AOzmond/usb-tree)](./LICENSE)
[![Release](https://img.shields.io/github/v/release/AOzmond/usb-tree)](https://github.com/AOzmond/usb-tree/releases)
[![Platform](https://img.shields.io/badge/platform-linux%20|%20windows-lightgrey.svg)](#)
[![Looking for Work](https://img.shields.io/badge/hiring-I'm%20looking%20for%20work-blue?style=flat-square)](https://aozmond.github.io)
[![Code style: gofumpt](https://img.shields.io/badge/code%20style-gofumpt-blue.svg)](https://github.com/mvdan/gofumpt)
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

## Overview

**USB Tree** is a cross-platform desktop and command-line tool that visualizes, logs, and monitors USB devices connected to your computer in real-time.  
It’s built in **Go** with **Wails v2**, combining a native Go backend with a modern **Svelte + TypeScript** frontend.

Use USB Tree to:
- View all currently connected USB devices in a hierarchical “tree” view  
- Log device connect/disconnect events with timestamps  
- Display vendor, and product information 
- Compare changes across sessions for debugging or audit purposes  

The tool is designed for developers, hardware engineers, and system administrators who want quick insight into USB topology and events.

---


## Features

- **Cross-Platform:** Runs on Linux and Windows  
- **Dual Interface:** CLI and GUI built from a shared Go library  
- **Real-Time Monitoring:** Detect hot-plug events and changes  
- **Device Metadata:** Vendor ID, Product ID, Path, and Bus info  
- **Lightweight & Fast:** Written in Go for efficiency and portability  

---

## Screenshots / Demo

TBD

---

## CLI

See [CLI documentation](cli/README.md) for usage details.

## GUI APP

See [GUI documentation](app/README.md) for usage details.

--- 
## License
This project is licensed under the GPL-2.0 License.

---
## Author

Alastair Ozmond  
Software Engineer | Full-Stack & Systems Developer  
[LinkedIn](www.linkedin.com/in/alastair-ozmond-108512179)
