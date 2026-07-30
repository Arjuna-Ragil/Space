package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type Window struct{
	Hwnd win.HWND
	Title string
}

func main() {
	hash := make(map[int][]Window)
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

		// Non essential program, remember to change as needed :)
		lowerTitle := strings.ToLower(title)
		if 	strings.Contains(lowerTitle, "windows input experience") ||
			strings.Contains(lowerTitle, "nvidia geforce overlay") ||
			strings.Contains(lowerTitle, "program manager") ||
			strings.Contains(lowerTitle, "yasbbar") ||
			strings.Contains(lowerTitle, "nahimic") {
			return 1
		}

		windowList = append(windowList, Window{Hwnd: hwnd, Title: title})

		hash[1] = windowList

		return 1
	})

	user32 := syscall.NewLazyDLL("user32.dll")
	enumWindows := user32.NewProc("EnumWindows")

	enumWindows.Call(enumFunc, 0)
	
	registerHotKey := user32.NewProc("RegisterHotKey")

	const (
		MOD_ALT = 0x0001
		WM_HOTKEY = 0x0312
	)

	// Basically alt + the workspace number
	registerHotKey.Call(0, 1, MOD_ALT, '1') 
	registerHotKey.Call(0, 2, MOD_ALT, '2')

	var msg win.MSG
	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		if ret == 0 {
			break
		}

		if msg.Message == WM_HOTKEY {
			pressedID := msg.WParam

			switch pressedID {
				case 1:
					fmt.Println("User pressed Alt+1! Switching to Workspace 1...")
				case 2:
					fmt.Println("User pressed Alt+2! Switching to Workspace 2...")
			}
		}

		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
