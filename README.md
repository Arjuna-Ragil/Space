# Space Workspace Manager

Space is a lightweight, lightning-fast virtual workspace and window manager for Windows. 

Written in Go, it operates entirely in the background from your System Tray, allowing you to seamlessly organize your windows across 9 virtual workspaces, instantly resize windows into predefined layouts, and manage your productivity without ever touching your mouse.

## Features

* **Virtual Workspaces:** Seamlessly switch between 9 distinct workspaces.
* **Instant Transitions:** Bypasses default Windows fade animations for immediate workspace switching.
* **Window Routing:** Instantly move your currently focused window to any other workspace.
* **Quick Resizing:** Snap windows to Full, Half, or Small layouts with customizable screen padding.
* **Dynamic Configuration:** Fully customizable hotkeys and padding via a simple `config.yaml` file.
* **System Tray Integration:** Clean, unobtrusive background execution.
* **Run on Startup:** Native toggle to start the application automatically when you log into Windows.
* **Graceful Shutdown:** Safely restores all hidden windows across all workspaces when you quit the application.

## Installation

1. Download the latest `Space.exe` release from the Releases tab.
2. Run the executable. It will appear in your System Tray.
3. On first launch, a `config.yaml` file is automatically generated in the same directory as the executable.

## Default Hotkeys

You can fully customize these hotkeys by right-clicking the System Tray icon and selecting **Open Config**.

### Workspace Navigation
* Switch to Workspace 1-9: `Alt + 1` through `Alt + 9`
* Move Active Window to Workspace 1-9: `Alt + Shift + 1` through `Alt + Shift + 9`

### Window Resizing
* Resize to Full (Padded): `Alt + Q`
* Resize to Half: `Alt + W`
* Resize to Small (Centered): `Alt + E`

## Configuration

The `config.yaml` file controls all keyboard shortcuts and layout padding. 

Example configuration:
```yaml
padding: 15
hotkeys:
  workspace1: Alt+1
  moveTo1: Alt+Shift+1
  resizeFull: Alt+Q
  resizeHalf: Alt+W
  resizeSmall: Alt+E
```

### Hotkey Modifiers
You can use any combination of the following modifiers separated by a `+` symbol:
* `Alt`
* `Shift`
* `Ctrl`
* `Win`

## Building from Source

To build Space from source, ensure you have Go 1.25 or later installed.

1. Clone the repository.
2. Install `go-winres` to embed the application icon and Windows manifest:
   
   ```bash
   go install github.com/tc-hib/go-winres@latest
   ```
   
3. Generate the Windows resources:
   
   ```bash
   go-winres make
   ```
   
4. Build the executable as a background GUI application (hides the console window):
   
   ```bash
   go build -ldflags "-H=windowsgui -s -w" -o Space.exe .
   ```

## System Window Filtering

Space intelligently ignores critical Windows system components to prevent accidental interference. Processes such as the Windows Input Experience, Program Manager (Desktop), and various overlays are safely excluded from workspace management. Task Manager is inherently ignored by Windows User Interface Privilege Isolation (UIPI) security.

# Thanks For Using Space :)
