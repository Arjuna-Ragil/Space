package main

import (
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type Window struct {
	Hwnd  win.HWND
	Title string
}

const (
	SPI_GETWORKAREA = 48
	SPI_SETWORKAREA = 47
)

const (
	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008
	WM_HOTKEY   = 0x0312
)

var currentHash = make(map[int][]Window)

func startHotkeyLoop() {
	runtime.LockOSThread()

	user32 := syscall.NewLazyDLL("user32.dll")
	registerHotKey := user32.NewProc("RegisterHotKey")

	reg := func(id int, hotkey string) {
		if hotkey == "" {
			return
		}
		mod, key := parseHotkey(hotkey)
		registerHotKey.Call(0, uintptr(id), mod, key)
	}

	// Register workspaces 1-9
	reg(1, currentConfig.Hotkeys.Workspace1)
	reg(2, currentConfig.Hotkeys.Workspace2)
	reg(3, currentConfig.Hotkeys.Workspace3)
	reg(4, currentConfig.Hotkeys.Workspace4)
	reg(5, currentConfig.Hotkeys.Workspace5)
	reg(6, currentConfig.Hotkeys.Workspace6)
	reg(7, currentConfig.Hotkeys.Workspace7)
	reg(8, currentConfig.Hotkeys.Workspace8)
	reg(9, currentConfig.Hotkeys.Workspace9)

	// Register move workspaces 10-18
	reg(10, currentConfig.Hotkeys.MoveTo1)
	reg(11, currentConfig.Hotkeys.MoveTo2)
	reg(12, currentConfig.Hotkeys.MoveTo3)
	reg(13, currentConfig.Hotkeys.MoveTo4)
	reg(14, currentConfig.Hotkeys.MoveTo5)
	reg(15, currentConfig.Hotkeys.MoveTo6)
	reg(16, currentConfig.Hotkeys.MoveTo7)
	reg(17, currentConfig.Hotkeys.MoveTo8)
	reg(18, currentConfig.Hotkeys.MoveTo9)

	// Register resizing scripts 19-21
	reg(19, currentConfig.Hotkeys.ResizeFull)
	reg(20, currentConfig.Hotkeys.ResizeHalf)
	reg(21, currentConfig.Hotkeys.ResizeSmall)

	var windowList []Window

	enumFunc := syscall.NewCallback(func(hwnd win.HWND, lParam uintptr) uintptr {
		if !win.IsWindowVisible(hwnd) {
			return 1
		}
		length := win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0)
		if length == 0 {
			return 1
		}
		buffer := make([]uint16, length+1)
		win.SendMessage(
			hwnd,
			win.WM_GETTEXT,
			uintptr(length+1),
			uintptr(unsafe.Pointer(&buffer[0])),
		)
		title := syscall.UTF16ToString(buffer)

		lowerTitle := strings.ToLower(title)
		if strings.Contains(lowerTitle, "windows input experience") ||
			strings.Contains(lowerTitle, "nvidia geforce overlay") ||
			strings.Contains(lowerTitle, "program manager") ||
			strings.Contains(lowerTitle, "yasbbar") ||
			strings.Contains(lowerTitle, "nahimic") {
			return 1
		}

		windowList = append(windowList, Window{Hwnd: hwnd, Title: title})
		return 1
	})

	enumWindows := user32.NewProc("EnumWindows")
	currentWorkspace := 1
	var msg win.MSG

	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		if ret == 0 {
			break
		}

		if msg.Message == WM_HOTKEY {
			pressedID := int(msg.WParam)

			if pressedID >= 19 && pressedID <= 21 {
				hwnd := win.GetForegroundWindow()
				if hwnd != 0 {
					win.ShowWindow(hwnd, win.SW_RESTORE)

					var workArea win.RECT
					win.SystemParametersInfo(SPI_GETWORKAREA, 0, unsafe.Pointer(&workArea), 0)

					pad := currentConfig.Padding

					var x, y, w, h int32
					fullW := workArea.Right - workArea.Left
					fullH := workArea.Bottom - workArea.Top

					switch pressedID {
					case 19:
						// Full padded
						w = fullW - (pad * 2)
						h = fullH - (pad * 2)
						x = workArea.Left + pad
						y = workArea.Top + pad
					case 20:
						// half
						w = (fullW - (pad * 3)) / 2
						h = fullH - (pad * 2)
						x = workArea.Left + pad
						y = workArea.Top + pad
					case 21:
						// small
						w = fullW / 2
						h = fullH / 2
						x = workArea.Left + (fullW - w) / 2
						y = workArea.Top + (fullH - h) / 2
					}

					win.SetWindowPos(hwnd, 0, x, y, w, h, win.SWP_NOZORDER|win.SWP_NOACTIVATE)
				}
				continue
			}

			targetWorkspace := pressedID
			isMove := false
			if pressedID >= 10 && pressedID <= 18 {
				targetWorkspace = pressedID - 9
				isMove = true
			}

			if targetWorkspace >= 1 && targetWorkspace <= 9 && targetWorkspace != currentWorkspace {
				if isMove {
					hwnd := win.GetForegroundWindow()
					if hwnd != 0 {
						length := win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0)
						title := ""
						if length > 0 {
							buffer := make([]uint16, length+1)
							win.SendMessage(hwnd, win.WM_GETTEXT, uintptr(length+1), uintptr(unsafe.Pointer(&buffer[0])))
							title = syscall.UTF16ToString(buffer)
						}

						lowerTitle := strings.ToLower(title)
						isSystemWindow := false
						if strings.Contains(lowerTitle, "windows input experience") ||
							strings.Contains(lowerTitle, "nvidia geforce overlay") ||
							strings.Contains(lowerTitle, "program manager") ||
							strings.Contains(lowerTitle, "yasbbar") ||
							strings.Contains(lowerTitle, "nahimic") || length == 0 {
							isSystemWindow = true
						}

						if !isSystemWindow {
							win.ShowWindow(hwnd, win.SW_HIDE)
							currentHash[targetWorkspace] = append(currentHash[targetWorkspace], Window{Hwnd: hwnd, Title: title})
						}
					}
				}

				windowList = nil
				enumWindows.Call(enumFunc, 0)

				currentWindows := make([]Window, len(windowList))
				copy(currentWindows, windowList)
				currentHash[currentWorkspace] = currentWindows

				for _, w := range currentHash[currentWorkspace] {
					win.ShowWindow(w.Hwnd, win.SW_HIDE)
				}

				for _, w := range currentHash[targetWorkspace] {
					win.ShowWindow(w.Hwnd, win.SW_SHOW)
					win.SetForegroundWindow(w.Hwnd)
				}

				currentWorkspace = targetWorkspace
			}
		}

		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
