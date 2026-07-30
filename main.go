package main

import (
	"fmt"
	"os"
	"os/signal"
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
	
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		for _, windows := range hash {
			for _, w := range windows {
				win.ShowWindow(w.Hwnd, win.SW_SHOW)
				win.SetForegroundWindow(w.Hwnd)
			}
		}
		os.Exit(0)
	}()

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
		return 1
	})

	user32 := syscall.NewLazyDLL("user32.dll")
	enumWindows := user32.NewProc("EnumWindows")
	
	registerHotKey := user32.NewProc("RegisterHotKey")

	const (
		MOD_ALT = 0x0001
		WM_HOTKEY = 0x0312
		MOD_SHIFT  = 0x0004
	)

	// The Hotkeys

	// Workspaces, basically alt + workspace number
	registerHotKey.Call(0, 1, MOD_ALT, '1') 
	registerHotKey.Call(0, 2, MOD_ALT, '2')
	registerHotKey.Call(0, 3, MOD_ALT, '3') 
	registerHotKey.Call(0, 4, MOD_ALT, '4')
	registerHotKey.Call(0, 5, MOD_ALT, '5') 
	registerHotKey.Call(0, 6, MOD_ALT, '6')
	registerHotKey.Call(0, 7, MOD_ALT, '7') 
	registerHotKey.Call(0, 8, MOD_ALT, '8')
	registerHotKey.Call(0, 9, MOD_ALT, '9')

	// Moving Program, basically alt + shift + workspace number
	registerHotKey.Call(0, 10, MOD_ALT | MOD_SHIFT, '1')
	registerHotKey.Call(0, 11, MOD_ALT | MOD_SHIFT, '2')
	registerHotKey.Call(0, 12, MOD_ALT | MOD_SHIFT, '3')
	registerHotKey.Call(0, 13, MOD_ALT | MOD_SHIFT, '4')
	registerHotKey.Call(0, 14, MOD_ALT | MOD_SHIFT, '5')
	registerHotKey.Call(0, 15, MOD_ALT | MOD_SHIFT, '6')
	registerHotKey.Call(0, 16, MOD_ALT | MOD_SHIFT, '7')
	registerHotKey.Call(0, 17, MOD_ALT | MOD_SHIFT, '8')
	registerHotKey.Call(0, 18, MOD_ALT | MOD_SHIFT, '9')

	currentWorkspace := 1
	var msg win.MSG
	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		if ret == 0 {
			break
		}

		if msg.Message == WM_HOTKEY {
			pressedID := int(msg.WParam)

			targetWorkspace := pressedID
			isMove := false
			if pressedID >= 10 && pressedID <= 18 {
				targetWorkspace = pressedID - 9
				isMove = true
			}

			if targetWorkspace >= 1 && targetWorkspace <= 9 && targetWorkspace != currentWorkspace {
				if isMove {
					fmt.Printf("Moving focused window to Workspace %d and switching...\n", targetWorkspace)
					
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
						if 	strings.Contains(lowerTitle, "windows input experience") ||
							strings.Contains(lowerTitle, "nvidia geforce overlay") ||
							strings.Contains(lowerTitle, "program manager") ||
							strings.Contains(lowerTitle, "yasbbar") ||
							strings.Contains(lowerTitle, "nahimic") || length == 0 {
							isSystemWindow = true
						}

						if !isSystemWindow {
							win.ShowWindow(hwnd, win.SW_HIDE)
							hash[targetWorkspace] = append(hash[targetWorkspace], Window{Hwnd: hwnd, Title: title})
						}
					}
				} else {
					fmt.Printf("Switching from Workspace %d to Workspace %d...\n", currentWorkspace, targetWorkspace)
				}

				windowList = nil
				enumWindows.Call(enumFunc, 0)
				
				currentWindows := make([]Window, len(windowList))
				copy(currentWindows, windowList)
				hash[currentWorkspace] = currentWindows

				for _, w := range hash[currentWorkspace] {
					win.ShowWindow(w.Hwnd, win.SW_HIDE)
				}

				for _, w := range hash[targetWorkspace] {
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
